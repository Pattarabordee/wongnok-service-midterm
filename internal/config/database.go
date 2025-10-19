package config

type Database struct {
	URL string `env:"DATABASE_URL" envDefault:"postgres://postgres:pass2word@localhost:5432/wongnok"`
}
