package global

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	level := logrus.WarnLevel
	if LogLevel >= 0 && LogLevel <= 6 {
		level = logrus.Level(LogLevel)
	}
	Log.SetLevel(level)

	if LogFile != "" {
		f, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			Log.WithField("func", "global.InitLogger").WithError(err).Warn("failed to open log file, logging to stdout only")
			return
		}
		Log.SetOutput(io.MultiWriter(os.Stdout, f))
	}
}
