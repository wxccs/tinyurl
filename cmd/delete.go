package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wxccs/tinyurl/app/Models"
	"github.com/wxccs/tinyurl/global"
	"github.com/wxccs/tinyurl/internal/config"
	"github.com/wxccs/tinyurl/internal/database"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <code>",
	Short: "Delete a short URL by its code",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		code := args[0]
		deleteShortUrl(code)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteShortUrl(code string) {
	global.InitLogger()
	funcName := "cmd.deleteShortUrl"

	cfg, err := config.Load()
	if err != nil {
		global.Log.WithField("func", funcName).WithError(err).Fatal("failed to load config")
	}
	global.Config = cfg

	db, err := database.Init(&cfg.Database)
	if err != nil {
		global.Log.WithField("func", funcName).WithError(err).Fatal("failed to init database")
	}
	global.DB = db
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	result := db.Where("short_code = ?", code).Delete(&Models.Url{})
	if result.Error != nil {
		global.Log.WithField("func", funcName).WithError(result.Error).Error("failed to delete short url")
		fmt.Fprintf(os.Stderr, "error: %v\n", result.Error)
		os.Exit(1)
	}

	if result.RowsAffected == 0 {
		global.Log.WithField("func", funcName).WithField("code", code).Warn("short code not found")
		fmt.Fprintf(os.Stderr, "short code '%s' not found\n", code)
		os.Exit(1)
	}

	global.Log.WithField("func", funcName).WithField("code", code).Info("short url deleted")
	fmt.Printf("short code '%s' deleted successfully\n", code)
}
