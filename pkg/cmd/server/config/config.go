package config

type Config struct {
	Port             string `env:"PORT" envDefault:"8080"`
	DatabaseKeyspace string `env:"KEYSPACE" envDefault:"moneway"`
}
