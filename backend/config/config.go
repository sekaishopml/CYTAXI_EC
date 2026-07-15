package config

type Config struct {
	App    AppConfig
	DB     DBConfig
	Redis  RedisConfig
	Log    LogConfig
	Auth   AuthConfig
	Otel   OtelConfig
}

type AppConfig struct {
	Env   string
	Name  string
	Port  int
	Debug bool
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type LogConfig struct {
	Level  string
	Format string
}

type AuthConfig struct {
	JWTSecret   string
	JWTExpiry   string
}

type OtelConfig struct {
	ServiceName string
	OTLPEndpoint string
}

type Loader interface {
	Load() (*Config, error)
}
