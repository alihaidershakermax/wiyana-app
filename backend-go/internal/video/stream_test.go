package video_test
import (
	"bytes"; "context"; "os"; "path/filepath"; "testing"; "time"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/store"
	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/video"
)
func TestGenerateEncrypted(t *testing.T) {
	s := video.NewStreamer("ffmpeg")
	cacheDir, _ := filepath.Abs("/tmp/wiyana-hls-cache/test-user/video-1")
	grant := store.Grant{Mode: "hls", Source: "/tmp/wiyana-uploads/test.mp4", WatermarkText: "test user", CacheDir: cacheDir}
	os.RemoveAll("/tmp/wiyana-hls-cache")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	playlist, err := s.Generate(ctx, grant)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	data, err := os.ReadFile(playlist)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("playlist:\n%s", string(data))
	key, err := os.ReadFile(video.KeyPath(grant.CacheDir))
	if err != nil {
		t.Fatalf("key missing: %v", err)
	}
	if len(key) != 16 {
		t.Fatalf("key len=%d want 16", len(key))
	}
	seg, err := os.ReadFile(filepath.Join(grant.CacheDir, "segment_00000.ts"))
	if err != nil {
		t.Fatal(err)
	}
	// An unencrypted TS file starts with 0x47 (MPEG-TS sync byte); an
	// encrypted segment must not.
	if bytes.HasPrefix(seg, []byte{0x47}) {
		t.Fatal("segment appears to be unencrypted")
	}
}
