package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env      string          `envconfig:"ENV" default:"local"`
	Postgres PostgresStorage `envconfig:"POSTGRES" required:"true"`
	HTTP     HTTPServer      `envconfig:"HTTP_SERVER" required:"true"`
	UrlInfo  string          `envconfig:"INFO_URL" required:"true"`
}

type PostgresStorage struct {
	StorageURL     string `envconfig:"STORAGE_URL" required:"true"`
	MigrationsPath string `envconfig:"MIGRATIONS_PATH" required:"true"`
}

type HTTPServer struct {
	Address     string        `envconfig:"ADDRESS" default:"localhost:8080"`
	Timeout     time.Duration `envconfig:"TIMEOUT" default:"4s"`
	IdleTimeout time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
	WithTimeout time.Duration `envconfig:"WITH_TIMEOUT" default:"10s"`
}

func MustLoad() *Config {
	var cfg Config

	//if err := godotenv.Load(".env"); err != nil {
	//	log.Println("Не удалось загрузить файл .env", err)
	//}

	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при парсинге конфигурации: %s", err)
	}
	return &cfg
}
