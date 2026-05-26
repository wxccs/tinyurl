package Routes

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxccs/tinyurl/app/Controllers/PageController"
	"github.com/wxccs/tinyurl/app/Controllers/UrlController"
	"github.com/wxccs/tinyurl/app/Middlewares"
)

func Setup(r *gin.Engine, staticFS fs.FS) {
	r.Use(Middlewares.CORS())

	PageController.SetStaticFS(staticFS)

	r.GET("/", PageController.ServeIndex)

	if staticFS != nil {
		assetsFS, err := fs.Sub(staticFS, "assets")
		if err == nil {
			r.StaticFS("/assets", http.FS(assetsFS))
		}

		// Register static files in dist root (e.g. favicon.svg, beian_logo.png)
		// before /:code so they take priority over the param route
		entries, err := fs.ReadDir(staticFS, ".")
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				filePath := name
				r.GET("/"+name, func(c *gin.Context) {
					data, err := fs.ReadFile(staticFS, filePath)
					if err != nil {
						c.String(http.StatusInternalServerError, "failed to read file")
						return
					}
					c.Data(http.StatusOK, mimeByExt(name), data)
				})
			}
		}

		r.NoRoute(func(c *gin.Context) {
			c.String(http.StatusNotFound, "not found")
		})
	}

	r.POST("/api/shorten", UrlController.Shorten)
	r.GET("/api/config", PageController.GetPublicConfig)
	r.GET("/:code", UrlController.Redirect)
}

func mimeByExt(name string) string {
	switch {
	case len(name) >= 4 && name[len(name)-4:] == ".png":
		return "image/png"
	case len(name) >= 4 && name[len(name)-4:] == ".svg":
		return "image/svg+xml"
	case len(name) >= 4 && name[len(name)-4:] == ".ico":
		return "image/x-icon"
	case len(name) >= 5 && name[len(name)-5:] == ".json":
		return "application/json"
	case len(name) >= 5 && name[len(name)-5:] == ".html":
		return "text/html; charset=utf-8"
	default:
		return "application/octet-stream"
	}
}
