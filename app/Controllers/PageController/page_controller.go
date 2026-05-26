package PageController

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxccs/tinyurl/app/Controllers"
	"github.com/wxccs/tinyurl/global"
)

type PublicConfig struct {
	Title string      `json:"title"`
	Beian BeianConfig `json:"beian"`
}

type BeianConfig struct {
	MIIT string `json:"miit"`
	MPS  string `json:"mps"`
}

var staticFS fs.FS

func SetStaticFS(fsys fs.FS) {
	staticFS = fsys
}

func ServeIndex(c *gin.Context) {
	funcName := "app.Controllers.PageController.ServeIndex"

	if staticFS == nil {
		global.Log.WithField("func", funcName).Warn("frontend not available")
		c.String(http.StatusNotFound, "frontend not available")
		return
	}

	data, err := fs.ReadFile(staticFS, "index.html")
	if err != nil {
		global.Log.WithField("func", funcName).WithError(err).Error("failed to read index.html")
		c.String(http.StatusInternalServerError, "failed to load page")
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

func GetPublicConfig(c *gin.Context) {
	funcName := "app.Controllers.PageController.GetPublicConfig"

	cfg := PublicConfig{
		Title: global.Config.Page.Title,
		Beian: BeianConfig{
			MIIT: global.Config.Beian.MIIT,
			MPS:  global.Config.Beian.MPS,
		},
	}

	global.Log.WithField("func", funcName).Debug("returning public config")
	Controllers.NewResponse(Controllers.RESPONSE_OK, "success", cfg).ResponseJson(200, c)
}
