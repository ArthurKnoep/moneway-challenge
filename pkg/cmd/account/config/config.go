package config

type Config struct {
	AccountPort      string `env:"ACCOUNT_PORT" envDefault:"8080"`
	DatabaseKeyspace string `env:"DB_KEYSPACE" envDefault:"moneway"`
	DatabaseHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
}
