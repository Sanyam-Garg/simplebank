package util

import "github.com/spf13/viper"

// Config stores all the configurations of the application
// Values are read by viper from a config file or env variables
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error){
	viper.AddConfigPath(path)
	viper.SetConfigName("app") // app.env
	viper.SetConfigType("env")

	// If the same key of a config property exists in the env variables as well, then load that one instead of the file one
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil{
		return
	}

	err = viper.Unmarshal(&config)
	return
}
