package config

type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:8080"`
	ServicePort      int    `env:"SERVICE_PORT"          envDefault:"8080"`
	ServiceHost      string `env:"SERVICE_HOST"          envDefault:"0.0.0.0"`
	PostgresUsername string `env:"POSTGRES_USER"         envDefault:"go-book"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"     envDefault:"go-book"`
	PostgresDB       string `env:"POSTGRES_DB"           envDefault:"go-book"`
	PostgresPort     int    `env:"POSTGRES_PORT"         envDefault:"5432"`
	MigrationsPath   string `env:"DB_MIGRATIONS_PATH"    envDefault:"file://migrations"`
	LogFilePath      string `env:"LOG_FILE_PATH"         envDefault:"logfile.log"`
	PostgresConn     string `env:"DB_CONNECTION_STRING"  envDefault:"postgres://go-book:go-book@go-book-postgres:5432/go-book?sslmode=disable"`
	UseFile          bool   `env:"USE_FILE"              envDefault:"false"`
}
