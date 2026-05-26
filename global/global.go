package global

import (
	"io/fs"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/wxccs/tinyurl/internal/config"
	"github.com/wxccs/tinyurl/internal/shortid"
)

var (
	Log      *logrus.Logger
	LogLevel int
	LogFile  string

	DB        *gorm.DB
	Config    *config.Config
	Generator *shortid.Generator
	StaticFS  fs.FS
)
