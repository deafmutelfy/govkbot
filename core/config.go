package core

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Token           string `yaml:"token" env:"TOKEN"`
	ConfirmationKey string `yaml:"confirmation_key" env:"CONFIRMATION_KEY"`
	Port            string `yaml:"port" env:"PORT"`
	RedisUrl        string `yaml:"redis_url" env:"REDIS_URL"`
	BotOwnerId      string `yaml:"bot_owner_id" env:"BOT_OWNER_ID"`
}

func (s *Config) Load(filename string) error {
	if err := cleanenv.ReadConfig(filename, s); err != nil {
		return err
	}

	return nil
}
