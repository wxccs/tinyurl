package main

import (
	"embed"
	"io/fs"

	"github.com/wxccs/tinyurl/cmd"
	"github.com/wxccs/tinyurl/global"
)

//go:embed web/dist
var embeddedFiles embed.FS

func init() {
	sub, err := fs.Sub(embeddedFiles, "web/dist")
	if err == nil {
		global.StaticFS = sub
	}
}

func main() {
	cmd.Execute()
}
