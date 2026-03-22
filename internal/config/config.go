package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBConfig    DBConfig          `yaml:"db"`
	Env         string            `yaml:"env" env-default:"local"`
	TGBotConfig TelegramBotConfig `yaml:"tg_bot_config"`
}
type DBConfig struct {
	Username      string        `yaml:"username"`
	Password      string        `yaml:"password"`
	Port          string        `yaml:"port"`
	Host          string        `yaml:"host"`
	Database      string        `yaml:"database"`
	MaxAttempts   int           `yaml:"max_attempts"`
	DelayAttempts time.Duration `yaml:"delay_attempts"`
}
type TelegramBotConfig struct {
	APIToken string        `yaml:"api_token"`
	PollTime time.Duration `yaml:"poll_time"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read congif : " + err.Error())
	}
	return &cfg
}
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
			res = envPath
		} else {
			res = "config/local.yaml" // default path for local development
		}
	}

	return res
}
