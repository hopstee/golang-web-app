package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Server struct {
	Addr    string `mapstructure:"addr"`
	BaseURL string `mapstructure:"base_url"`
}

type Database struct {
	Driver     string `mapstructure:"driver"`
	DataSource string `mapstructure:"dsn"`
}

type KVStore struct {
	Type  string `mapstructure:"type"`
	Redis Redis  `mapstructure:"redis"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type FileStorage struct {
	Type  string       `mapstructure:"type"`
	Local LocalStorage `mapstructure:"local"`
}

type LocalStorage struct {
	BasePath string `mapstructure:"base_path"`
}

type JWTAuth struct {
	Secret string `mapstructure:"secret"`
}

type Authentication struct {
	JWT JWTAuth `mapstructure:"jwt"`
}

type Telegram struct {
	Token  string `mapstructure:"token"`
	ChatID string `mapstructure:"chat_id"`
}

type Static struct {
	Dir string `mapstructure:"dir"`
}

type Webapp struct {
	Dir string `mapstructure:"dir"`
}

type Schemas struct {
	Schema  string `mapstructure:"schema"`
	Pages   string `mapstructure:"pages"`
	Layouts string `mapstructure:"layouts"`
	Blocks  string `mapstructure:"blocks"`
	Modules string `mapstructure:"modules"`
	Shared  string `mapstructure:"shared"`
}

type Config struct {
	Server         Server         `mapstructure:"server"`
	Database       Database       `mapstructure:"database"`
	KVStore        KVStore        `mapstructure:"kvstore"`
	FileStorage    FileStorage    `mapstructure:"filestorage"`
	Authentication Authentication `mapstructure:"authentication"`
	Telegram       Telegram       `mapstructure:"telegram"`
	Static         Static         `mapstructure:"static"`
	Webapp         Webapp         `mapstructure:"webapp"`
	Schemas        Schemas        `mapstructure:"schemas"`
}

func LoadConfig(path string) (*Config, error) {
	_ = godotenv.Load()

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yml")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	for _, key := range v.AllKeys() {
		val := v.GetString(key)
		if strings.Contains(val, "${") {
			v.Set(key, os.ExpandEnv(val))
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
