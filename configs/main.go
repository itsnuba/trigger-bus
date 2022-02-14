package configs

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	configDb `mapstructure:",squash"`

	Debug bool `mapstructure:"debug"`
}

func Load() *Config {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))

	viper.SetConfigFile(env + ".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err.Error())
	}

	return &config
}
