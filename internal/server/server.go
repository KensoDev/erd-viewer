package server

// This package is deprecated. Use github.com/kensodev/erd-viewer/pkg/webview instead.
// These type aliases and wrappers are provided for backward compatibility.

import (
	"github.com/kensodev/erd-viewer/internal/db"
	"github.com/kensodev/erd-viewer/pkg/webview"
	"github.com/kensodev/erd-viewer/web"
)

// Deprecated: Use webview.Server instead
type Server = webview.Server

// Deprecated: Use webview.New instead
func New(data *db.SchemaData, listenAddr string) (*Server, error) {
	return webview.New(webview.Config{
		SchemaData: data,
		ListenAddr: listenAddr,
		Assets:     &webview.EmbedAssets{FS: web.Files},
	})
}
