package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv"
)

type Config struct {
	Env     string   `env:"APP_ENV" env-default:"local" env-required:"true"`
	DB      DBConfig `env-prefix:"DB_"`
	GRPCApp GRPCApp  `env-prefix:"APP_"`
}

type DBConfig struct {
	Host     string `env:"HOST" env-required:"true"`
	Port     int    `env:"PORT" env-default:"5432"`
	User     string `env:"USER" env-default:"postgres"`
	Password string `env:"PASSWORD" env-required:"true"`
	Name     string `env:"NAME" env-required:"true"`
	SSLMode  string `env:"SSLMode" env-default:"disable"`
}

func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode,
	)
}

type GRPCApp struct {
	Port    int64 `env:"PORT" env-default:"localhost:44044"`
	Timeout int64 `env:"TIMEOUT" env-default:"4000000000"`
}

func (g GRPCApp) AsDuration() time.Duration {
	return time.Duration(g.Timeout)
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Error("failed to load config from env")
	}

	if cfg.Env == "" {
		slog.Error("APP_ENV cannot be empty")
	}
	if cfg.DB.DSN() == "" {
		slog.Error("database configuration is incomplete")
	}

	slog.Info(fmt.Sprintf("config loaded: env=%s, db=%s, port=%d",
		cfg.Env, cfg.DB.Host, cfg.GRPCApp.Port))

	return &cfg
}
