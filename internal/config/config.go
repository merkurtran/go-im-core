package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig   `mapstructure:"app"`
	MongoDB MongoConfig `mapstructure:"mongo"`
	JWT     JWTConfig   `mapstructure:"jwt"`
	HTTP    HTTPConfig  `mapstructure:"http"`
	Log     LogConfig   `mapstructure:"log"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Env         string `mapstructure:"env"` // dev, prod
	Description string `mapstructure:"description"`
}

type MongoConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Timeout  int    `mapstructure:"timeout"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

type HTTPConfig struct {
	Port    int `mapstructure:"port"`
	Timeout int `mapstructure:"timeout"`
}

type LogConfig struct {
	Level   string `mapstructure:"level"`
	File    string `mapstructure:"file"`
	MaxSize int    `mapstructure:"max_size"` // MB
	MaxAge  int    `mapstructure:"max_age"`  // days
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath("./")

	setDefaults(v)

	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config.yaml: %v", err)
	}

	env := getEnv()
	envConfigName := "config." + env

	v.SetConfigName(envConfigName)
	if err := v.MergeInConfig(); err == nil {
		log.Printf("Merged env config: %s.yaml\n", envConfigName)
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		log.Printf("Warning: failed to read %s.yaml: %v\n", envConfigName, err)
	}

	v.SetEnvPrefix("CONFIG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	return &config, nil
}

// getEnv 获取当前运行环境，默认 dev
func getEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	return env
}

func (c *Config) Validate() error {
	if c.App.Env == "" {
		return errors.New("app.env is required")
	}
	if c.MongoDB.URI == "" {
		return errors.New("mongo.uri is required")
	}
	if c.MongoDB.Database == "" {
		return errors.New("mongo.database is required")
	}
	if c.MongoDB.Username == "" {
		return errors.New("mongo.username is required")
	}
	if c.MongoDB.Password == "" {
		return errors.New("mongo.password is required")
	}
	if c.MongoDB.Timeout == 0 {
		return errors.New("mongo.timeout is required")
	}
	if c.JWT.Secret == "" {
		return errors.New("jwt.secret is required")
	}
	if c.JWT.Expire == 0 {
		return errors.New("jwt.expire is required")
	}
	if c.HTTP.Port < 1 || c.HTTP.Port > 65535 {
		return errors.New("http.port is invalid")
	}
	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "go-im-core")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.env", "dev")
	v.SetDefault("http.port", 8080)
	v.SetDefault("http.timeout", 30)
	v.SetDefault("log.level", "info")
	v.SetDefault("mongo.timeout", 10)
}
