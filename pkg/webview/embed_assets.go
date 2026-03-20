package webview

import "io/fs"

// EmbedAssets provides assets from an embed.FS
type EmbedAssets struct {
	FS fs.FS
}

// ReadFile reads a file from the embedded filesystem
func (e *EmbedAssets) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(e.FS, name)
}
