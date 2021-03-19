package config

type APIConfig struct {
	ListenAddress string `split_words:"true" default:"0.0.0.0:80"`
}

type PostgresConfig struct {
	Host     string `required:"true"`
	Port     int    `default:"5432"`
	User     string `required:"true"`
	Password string `required:"true"`
	Database string `required:"true"`
	SSLMode  string `split_words:"true" default:"disable"`
	Debug    bool   `default:"false"`
}

type MigrationConfig struct {
	Postgres PostgresConfig `split_words:"true"`
}
