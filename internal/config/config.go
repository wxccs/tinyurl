package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	ShortURL ShortURLConfig `mapstructure:"shorturl"`
	Page     PageConfig     `mapstructure:"page"`
	Beian    BeianConfig    `mapstructure:"beian"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Path     string `mapstructure:"path"`
}

type ShortURLConfig struct {
	Length int `mapstructure:"length"`
	NodeID int `mapstructure:"node_id"`
}

type PageConfig struct {
	Title string `mapstructure:"title"`
}

type BeianConfig struct {
	MIIT string `mapstructure:"miit"`
	MPS  string `mapstructure:"mps"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	cfg.validate()

	return cfg, nil
}

func (c *Config) validate() {
	if c.Database.Port == 0 {
		switch c.Database.Type {
		case "mysql":
			c.Database.Port = 3306
		case "postgres":
			c.Database.Port = 5432
		}
	}
	if c.ShortURL.Length < 7 {
		c.ShortURL.Length = 7
	}
	if c.ShortURL.NodeID < 0 {
		c.ShortURL.NodeID = 0
	}
	if c.ShortURL.NodeID > 15 {
		c.ShortURL.NodeID = 15
	}
}
