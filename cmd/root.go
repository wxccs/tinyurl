package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wxccs/tinyurl/app/Routes"
	"github.com/wxccs/tinyurl/global"
	"github.com/wxccs/tinyurl/internal/config"
	"github.com/wxccs/tinyurl/internal/database"
	"github.com/wxccs/tinyurl/internal/shortid"
)

var cfgFile string
var logLevel int
var logFile string

var rootCmd = &cobra.Command{
	Use:   "tinyurl",
	Short: "A short URL generator service",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
	Version: global.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/tinyurl/config.yaml)")
	rootCmd.PersistentFlags().IntVar(&logLevel, "log-level", 3, "log level (0=Panic, 1=Fatal, 2=Error, 3=Warn, 4=Info, 5=Debug, 6=Trace)")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "log file path (default empty, logs to stdout only)")
	_ = viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("log.file", rootCmd.PersistentFlags().Lookup("log-file"))
	rootCmd.SetVersionTemplate(fmt.Sprintf(`{{with .Name}}{{printf "%%s version information: " .}}{{end}}
   {{printf "Version:    %%s" .Version}}
   Build Time:          %s
   Git Revision:        %s
   Go version:          %s
   OS/Arch:                     %s/%s
`, global.BuildTime, global.GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(fmt.Sprintf("%s%c.config%ctinyurl", home, os.PathSeparator, os.PathSeparator))
		viper.AddConfigPath("/etc/tinyurl")
		viper.AddConfigPath("/usr/etc/tinyurl")
		viper.AddConfigPath("/usr/local/etc/tinyurl")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("TU")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func runServer() {
	global.LogLevel = logLevel
	global.LogFile = logFile
	global.InitLogger()
	global.Log.WithField("func", "cmd.runServer").Info("initializing tinyurl service")

	cfg, err := config.Load()
	if err != nil {
		global.Log.WithField("func", "cmd.runServer").WithError(err).Fatal("failed to load config")
	}
	global.Config = cfg

	db, err := database.Init(&cfg.Database)
	if err != nil {
		global.Log.WithField("func", "cmd.runServer").WithError(err).Fatal("failed to init database")
	}
	global.DB = db
	global.Log.WithField("func", "cmd.runServer").WithField("type", cfg.Database.Type).Info("database initialized")

	global.Generator = shortid.NewGenerator(cfg.ShortURL.NodeID, cfg.ShortURL.Length)
	global.Log.WithField("func", "cmd.runServer").WithField("length", cfg.ShortURL.Length).Info("short id generator initialized")

	if global.LogLevel < 5 {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	Routes.Setup(r, global.StaticFS)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		global.Log.WithField("func", "cmd.runServer").WithField("addr", addr).Info("server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.WithField("func", "cmd.runServer").WithError(err).Fatal("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	global.Log.WithField("func", "cmd.runServer").Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		global.Log.WithField("func", "cmd.runServer").WithError(err).Error("server forced to shutdown")
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	global.Log.WithField("func", "cmd.runServer").Info("server exited")
}
