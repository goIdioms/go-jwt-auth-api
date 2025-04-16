package database

import (
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	DBUri               string        `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri            string        `mapstructure:"REDIS_URL"`
	Port                string        `mapstructure:"PORT"`
	AccessJwtSecret     string        `mapstructure:"ACCESS_JWT_SECRET"`
	AccessJwtExpiresIn  time.Duration `mapstructure:"ACCESS_JWT_EXPIRED_IN"`
	AccessJwtMaxAge     int           `mapstructure:"ACCESS_JWT_MAXAGE"`
	RefreshJwtSecret    string        `mapstructure:"REFRESH_JWT_SECRET"`
	RefreshJwtExpiresIn time.Duration `mapstructure:"REFRESH_JWT_EXPIRED_IN"`
	RefreshJwtMaxAge    int           `mapstructure:"REFRESH_JWT_MAXAGE"`
	ClientOrigin        string        `mapstructure:"CLIENT_ORIGIN"`
}

var UserCollection *mongo.Collection

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
