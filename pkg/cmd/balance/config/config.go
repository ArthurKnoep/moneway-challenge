package config

type Config struct {
	BalancePort      string `env:"BALANCE_PORT" envDefault:"8081"`
	DatabaseKeyspace string `env:"DB_KEYSPACE" envDefault:"moneway"`
	DatabaseHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
}
