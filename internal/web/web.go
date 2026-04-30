package web

import "embed"

// FS contains the embedded static web assets served by the HTTP server.
//
//go:embed index.html
var FS embed.FS
