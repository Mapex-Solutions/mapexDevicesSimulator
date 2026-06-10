package web

import (
	"embed"
	"io/fs"
)

// distFS holds the built SPA, copied into dist/ by the production build before
// `go build` so the sidecar binary stays self-contained. In dev the dist/ holds
// only a placeholder and the Vite server serves the real SPA.
//
//go:embed all:dist
var distFS embed.FS

// SPA returns the embedded SPA rooted at dist/, or nil when the subtree is
// missing.
func SPA() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil
	}
	return sub
}
