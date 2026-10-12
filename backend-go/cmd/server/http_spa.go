package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// newSPAServer returns a handler that serves real files from the SPA build
// directory and falls back to index.html for every other GET path, so
// front-end routers like /admin keep working when the app is reloaded.
func newSPAServer(root string) http.Handler {
	fs := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			writeError(w, http.StatusMethodNotAllowed, "طريقة الطلب غير مدعومة")
			return
		}
		clean := strings.TrimPrefix(filepath.ToSlash(r.URL.Path), "/")
		target := filepath.Join(root, filepath.FromSlash(clean))
		if info, err := os.Stat(target); err == nil && !info.IsDir() {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			fs.ServeHTTP(w, r)
			return
		}
		// Serve index.html for every virtual route (/admin, /courses/..., ...).
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFile(w, r, filepath.Join(root, "index.html"))
	})
}
