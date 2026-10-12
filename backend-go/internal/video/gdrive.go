package video

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// GoogleDriveDownloader fetches a video hosted on Google Drive (private or
// public) and saves it to a local file so the normal encrypted-HLS pipeline
// can process it. Use a service-account JSON (GOOGLE_DRIVE_SA_JSON) for
// private files; otherwise an API key with shared-link access is tried.
type GoogleDriveDownloader struct {
	client    *http.Client
	apiKey    string
	saTokenFn func(ctx context.Context) (string, error)
}

type gdriveServiceAccount struct {
	TokenURI      string `json:"token_uri"`
	ClientEmail   string `json:"client_email"`
	PrivateKeyID  string `json:"private_key_id"`
	PrivateKey    string `json:"private_key"`
	ProjectID     string `json:"project_id"`
	Type          string `json:"type"`
}

func NewGoogleDriveDownloader(serviceAccountJSON, apiKey string) (*GoogleDriveDownloader, error) {
	serviceAccountJSON = strings.TrimSpace(serviceAccountJSON)
	apiKey = strings.TrimSpace(apiKey)
	d := &GoogleDriveDownloader{
		client: &http.Client{Timeout: 5 * time.Minute},
		apiKey: apiKey,
	}
	if serviceAccountJSON != "" {
		var sa gdriveServiceAccount
		if err := json.Unmarshal([]byte(serviceAccountJSON), &sa); err != nil {
			return nil, fmt.Errorf("parse service account JSON: %w", err)
		}
		if sa.PrivateKey == "" || sa.ClientEmail == "" {
			return nil, errors.New("service account JSON must include private_key and client_email")
		}
		d.saTokenFn = func(ctx context.Context) (string, error) {
			return saSignedJWTToken(ctx, sa)
		}
	}
	if d.saTokenFn == nil && d.apiKey == "" {
		return nil, errors.New("google drive support requires GOOGLE_DRIVE_SA_JSON or GOOGLE_DRIVE_API_KEY")
	}
	return d, nil
}

func (d *GoogleDriveDownloader) Enabled() bool { return d.saTokenFn != nil || d.apiKey != "" }

// saSignedJWTToken builds and signs an RS256 JWT for the Drive API, then
// exchanges it for a short-lived access token.
func saSignedJWTToken(ctx context.Context, sa gdriveServiceAccount) (string, error) {
	block, _ := pem.Decode([]byte(sa.PrivateKey))
	if block == nil {
		return "", errors.New("invalid private key PEM")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", errors.New("unsupported private key format")
		}
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("private key is not RSA")
	}

	now := time.Now()
	claims := map[string]any{
		"iss":   sa.ClientEmail,
		"scope": "https://www.googleapis.com/auth/drive.readonly",
		"aud":   sa.TokenURI,
		"exp":   now.Add(time.Hour).Unix(),
		"iat":   now.Unix(),
		}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString(claimsBytes)
	sigInput := header + "." + payload
	digest := sha256.Sum256([]byte(sigInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, digest[:])
	if err != nil {
		return "", err
	}
	jwt := sigInput + "." + base64.RawURLEncoding.EncodeToString(sig)

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, sa.TokenURI, strings.NewReader(url.Values{
		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  {jwt},
	}.Encode()))
	if reqErr != nil {
		return "", reqErr
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, respErr := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if respErr != nil {
		return "", respErr
	}
	defer resp.Body.Close()
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return "", err
	}
	if tok.AccessToken == "" {
		return "", errors.New("service account token exchange failed")
	}
	return tok.AccessToken, nil
}

// ParseID extracts the Google Drive file id from various URL formats.
// Supported: gdrive://{fileId}, full share links, /d/{fileId} paths.
func ParseID(source string) string {
	s := strings.TrimSpace(source)
	if strings.HasPrefix(s, "gdrive://") {
		return strings.TrimPrefix(s, "gdrive://")
	}
	if u, err := url.Parse(s); err == nil && u.Hostname() != "" {
		for _, pat := range []string{"/file/d/", "/d/"} {
			if i := strings.Index(s, pat); i >= 0 {
				rest := s[i+len(pat):]
				if j := strings.IndexAny(rest, "/"); j > 0 {
					rest = rest[:j]
				}
				return rest
			}
		}
		if q := u.Query().Get("id"); q != "" {
			return q
		}
	}
	if len(s) >= 10 && !strings.ContainsAny(s, "/:") {
		return s // raw file id
	}
	return ""
}

// Download fetches the Drive file into dst. Returns an error if the file is
// not accessible (quota exceeded or access denied).
func (d *GoogleDriveDownloader) Download(ctx context.Context, fileID, dst string) error {
	var body io.ReadCloser
	var err error

	if d.saTokenFn != nil {
		token, terr := d.saTokenFn(ctx)
		if terr == nil {
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
				"https://www.googleapis.com/drive/v3/files/"+url.PathEscape(fileID)+"?alt=media&supportsAllDrives=true", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, respErr := d.client.Do(req)
			if respErr == nil {
				if resp.StatusCode == http.StatusOK {
					body = resp.Body
				} else {
					resp.Body.Close()
					err = fmt.Errorf("service account fetch: HTTP %d", resp.StatusCode)
				}
			} else {
				err = respErr
			}
		}
	}

	if body == nil && d.apiKey != "" {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
			"https://www.googleapis.com/drive/v3/files/"+url.PathEscape(fileID)+"?alt=media&supportsAllDrives=true&key="+url.QueryEscape(d.apiKey), nil)
		resp, respErr := d.client.Do(req)
		if respErr == nil {
			if resp.StatusCode == http.StatusOK {
				body = resp.Body
			} else {
				resp.Body.Close()
				if err == nil {
					err = fmt.Errorf("api key fetch: HTTP %d", resp.StatusCode)
				}
			}
		} else if err == nil {
			err = respErr
		}
	}

	if body == nil {
		return fmt.Errorf("google drive file %s not accessible: %w", fileID, err)
	}
	defer body.Close()

	out, cerr := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if cerr != nil {
		return cerr
	}
	defer out.Close()
	_, cerr = io.Copy(out, body)
	if cerr != nil {
		return cerr
	}
	return nil
}
