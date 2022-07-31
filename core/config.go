package core

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Token           string `yaml:"token" env:"TOKEN"`
	UserToken       string `yaml:"user_token" env:"USER_TOKEN"`
	ConfirmationKey string `yaml:"confirmation_key" env:"CONFIRMATION_KEY"`
	Port            string `yaml:"port" env:"PORT"`
	RedisUrl        string `yaml:"redis_url" env:"REDIS_URL"`
	BotOwnerId      string `yaml:"bot_owner_id" env:"BOT_OWNER_ID"`
	Host            string `yaml:"host" env:"HOST"`
	EventManager    string `yaml:"event_manager" env:"EVENT_MANAGER"`

	GroupId int
}

func (s *Config) Load(filename string) error {
	if err := cleanenv.ReadConfig(filename, s); err != nil {
		if err := cleanenv.ReadEnv(s); err != nil {
			return err
		}
	}

	return nil
}
