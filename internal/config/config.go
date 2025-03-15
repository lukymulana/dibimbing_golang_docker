package config

type Config struct {
	Port string `env:"PORT" envDefault:"8000"`
	Env  string `env:"ENV" envDefault:"dev"`
}
