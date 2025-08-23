package config

type Config struct {
	App    App           `yaml:"app" validate:"required"`
	DB     DB            `yaml:"db" validate:"required"`
	Redis  Redis         `yaml:"redis" validate:"required"`
	Secret SecreteConfig `yaml:"secret" validate:"required"`
}

type App struct {
	Port string `yaml:"port" validate:"required"`
}

type DB struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
}

type SecreteConfig struct {
	JWTSecret string `yaml:"jwt_secret" validate:"required"`
}

type Redis struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	Password string `yaml:"password"`
}
