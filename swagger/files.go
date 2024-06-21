package swagger

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var dist embed.FS

var FS, _ = fs.Sub(dist, "static")
