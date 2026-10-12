package video

import (
	"context"
	"crypto/aes"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/store"
)

type Streamer struct {
	ffmpeg string
	mu     sync.Mutex
	jobs   map[string]struct{}
}

func NewStreamer(ffmpeg string) *Streamer {
	return &Streamer{ffmpeg: ffmpeg, jobs: map[string]struct{}{}}
}

// KeyPath returns the path of the per-grant AES key file inside the grant cache.
// The segment URLs in the playlist point to this file and the backend serves it
// only to the same grant (device-bound, signed URL).
func KeyPath(cacheDir string) string { return filepath.Join(cacheDir, "stream.key") }

func (s *Streamer) Generate(ctx context.Context, grant store.Grant) (string, error) {
	if grant.Mode != "hls" {
		return "", errors.New("grant is not an HLS grant")
	}
	if grant.Source == "" {
		return "", errors.New("HLS source is empty")
	}
	source := grant.Source
	var tmpFile string
	if strings.HasPrefix(source, "gdrive://") {
		id := ParseID(source)
		if id == "" {
			return "", fmt.Errorf("invalid google drive source: %s", source)
		}
		dl, err := NewGoogleDriveDownloader(grant.GoogleDriveSA, grant.GoogleDriveAPIKey)
		if err != nil {
			return "", err
		}
		tmpFile = filepath.Join(os.TempDir(), fmt.Sprintf("gdrive-%s-%d.mp4", id[:8], time.Now().UnixNano()))
		defer os.Remove(tmpFile)
		if err := dl.Download(ctx, id, tmpFile); err != nil {
			return "", err
		}
		source = tmpFile
	}
	if _, err := os.Stat(source); err != nil {
		return "", fmt.Errorf("video source unavailable: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(grant.CacheDir), 0750); err != nil {
		return "", err
	}
	if err := os.MkdirAll(grant.CacheDir, 0750); err != nil {
		return "", err
	}
	playlist := filepath.Join(grant.CacheDir, "playlist.m3u8")
	ready := filepath.Join(grant.CacheDir, ".ready")
	if readyFile(playlist, ready) {
		return playlist, nil
	}

	s.mu.Lock()
	if _, running := s.jobs[grant.CacheDir]; running {
		s.mu.Unlock()
		return waitForReady(ctx, playlist, ready)
	}
	s.jobs[grant.CacheDir] = struct{}{}
	s.mu.Unlock()
	defer func() { s.mu.Lock(); delete(s.jobs, grant.CacheDir); s.mu.Unlock() }()
	if readyFile(playlist, ready) {
		return playlist, nil
	}

	lockDir := grant.CacheDir + ".lock"
	if err := os.Mkdir(lockDir, 0700); err != nil {
		if os.IsExist(err) {
			return waitForReady(ctx, playlist, ready)
		}
		return "", err
	}
	defer os.RemoveAll(lockDir)
	if readyFile(playlist, ready) {
		return playlist, nil
	}
	if _, err := exec.LookPath(s.ffmpeg); err != nil {
		return "", fmt.Errorf("ffmpeg is not installed: %w", err)
	}

	// Per-session 128-bit encryption key and IV. Each viewer gets a unique key
	// so a leaked link is useless without the same grant token.
	key, iv, err := newKeyPair()
	if err != nil {
		return "", err
	}
	// The key must live inside the final cache directory (so it stays next to
	// the segments it decrypts) but CacheDir is wiped before the atomic rename
	// below; write it into buildDir and move it into place afterwards. The
	// keyFile path is still needed by writeKeyInfo for the ffmpeg key info file.
	keyFile := KeyPath(grant.CacheDir)

	buildDir := fmt.Sprintf("%s.tmp-%d", grant.CacheDir, time.Now().UnixNano())
	if err := os.MkdirAll(buildDir, 0750); err != nil {
		return "", err
	}
	defer os.RemoveAll(buildDir)
	// Write the key into buildDir (it survives the atomic rename below) but
	// register its absolute final path so ffmpeg's key-info file can open it
	// reliably from any working directory.
	buildKeyAbs := filepath.Join(buildDir, "stream.key")
	if err := os.WriteFile(buildKeyAbs, key, 0600); err != nil {
		return "", err
	}
	if abs, err := filepath.Abs(buildKeyAbs); err == nil {
		buildKeyAbs = abs
	}
	buildPlaylist := filepath.Join(buildDir, "playlist.m3u8")
	segmentPattern := filepath.Join(buildDir, "segment_%05d.ts")
	args := []string{"-hide_banner", "-loglevel", "error", "-y", "-i", grant.Source}
	if grant.WatermarkText != "" {
		args = append(args, "-vf", "drawtext=fontcolor=white:fontsize=18:box=1:boxcolor=black@0.45:x=20:y=20:text='"+escapeFilter(grant.WatermarkText)+"'")
	}
	args = append(args,
		"-c:v", "libx264", "-preset", "veryfast",
		"-c:a", "aac", "-b:a", "128k",
		"-hls_time", "6", "-hls_playlist_type", "vod",
		"-hls_segment_filename", segmentPattern,
		"-hls_key_info_file", writeKeyInfo(buildDir, buildKeyAbs, iv),
		"-f", "hls", buildPlaylist)
	command := exec.CommandContext(ctx, s.ffmpeg, args...)
	if output, err := command.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg HLS generation failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	if !readyFile(buildPlaylist, "") {
		return "", errors.New("ffmpeg did not create a valid HLS playlist")
	}
	if err := os.WriteFile(filepath.Join(buildDir, ".ready"), []byte("ok\n"), 0600); err != nil {
		return "", err
	}
	_ = os.RemoveAll(grant.CacheDir)
	if err := os.Rename(buildDir, grant.CacheDir); err != nil {
		return "", err
	}
	// The key was built inside the temporary dir; move it into the final cache
	// alongside the playlist and segments it decrypts.
	if err := os.Rename(filepath.Join(grant.CacheDir, "stream.key"), keyFile); err != nil {
		return "", err
	}
	return playlist, nil
}

// writeKeyInfo writes an ffmpeg key-info file pointing at the grant key with its IV.
func writeKeyInfo(dir, keyFile string, iv []byte) string {
	path := filepath.Join(dir, "keyinfo.txt")
	_ = os.WriteFile(path, []byte(fmt.Sprintf("%s\n%s\n%s\n", keyFile, keyFile, hexstr(iv))), 0600)
	return path
}

func hexstr(b []byte) string {
	const h = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = h[v>>4]
		out[i*2+1] = h[v&0x0f]
	}
	return string(out)
}

func newKeyPair() (key, iv []byte, err error) {
	key = make([]byte, 16)
	if _, err = rand.Read(key); err != nil {
		return nil, nil, err
	}
	iv = make([]byte, aes.BlockSize)
	if _, err = rand.Read(iv); err != nil {
		return nil, nil, err
	}
	// Ensure the IV, when interpreted as a big-endian integer, is at least 1
	// so the AES-CTR counter never collides with the key stream.
	if isZero(iv) {
		iv[len(iv)-1] = 1
	}
	return key, iv, nil
}

func isZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

func readyFile(playlist, marker string) bool {
	playlistInfo, err := os.Stat(playlist)
	if err != nil || playlistInfo.Size() == 0 {
		return false
	}
	if marker == "" {
		return true
	}
	markerInfo, err := os.Stat(marker)
	return err == nil && !markerInfo.IsDir() && markerInfo.Size() > 0
}

func waitForReady(ctx context.Context, playlist, marker string) (string, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		if readyFile(playlist, marker) {
			return playlist, nil
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
		}
	}
}

func escapeFilter(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, "'", "\\'")
	text = strings.ReplaceAll(text, ":", "\\:")
	text = strings.ReplaceAll(text, "%", "\\%")
	return text
}
