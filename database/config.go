package database

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBUri        string        `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri     string        `mapstructure:"REDIS_URL"`
	Port         string        `mapstructure:"PORT"`
	JwtSecret    string        `mapstructure:"JWT_SECRET"`
	JwtExpiresIn time.Duration `mapstructure:"JWT_EXPIRED_IN"`
	JwtMaxAge    int           `mapstructure:"JWT_MAXAGE"`
	ClientOrigin string        `mapstructure:"CLIENT_ORIGIN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
