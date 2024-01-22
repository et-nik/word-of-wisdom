package config

import "github.com/caarlos0/env/v10"

type Config struct {
	ListenAddr string `env:"SERVER_ADDR" envDefault:""`
	ListenPort int    `env:"SERVER_PORT" envDefault:"9100"`

	DifficultyWidth  int `env:"DIFFICULTY_WIDTH" envDefault:"70"`
	DifficultyLength int `env:"DIFFICULTY_LENGTH" envDefault:"3"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

func validate(cfg *Config) error {
	if cfg.DifficultyWidth < 1 {
		return ValidationError("DIFFICULTY_WIDTH must be greater than 0")
	}

	if cfg.DifficultyWidth > 200 {
		return ValidationError("DIFFICULTY_WIDTH must be less than 200")
	}

	if cfg.DifficultyLength < 1 {
		return ValidationError("DIFFICULTY_LENGTH must be greater than 0")
	}

	if cfg.DifficultyLength > 5 {
		return ValidationError("DIFFICULTY_LENGTH must be less than 5")
	}

	return nil
}
