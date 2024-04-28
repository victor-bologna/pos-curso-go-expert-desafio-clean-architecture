package configs

import (
	"github.com/spf13/viper"
)

var cfg *Conf

type Conf struct {
	DB  `mapstructure:"DB"`
	API `mapstructure:"API"`
}

type DB struct {
	Driver   string `mapstructure:"DRIVER"`
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	Name     string `mapstructure:"NAME"`
}

type API struct {
	WebServerPort     string `mapstructure:"PORT"`
	GRPCServerPort    string `mapstructure:"GRPC"`
	GraphQLServerPort string `mapstructure:"GRAPHQL"`
	Rabbit            string `mapstructure:"RABBIT"`
}

func LoadConfig(path string) *Conf {
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	return cfg
}
