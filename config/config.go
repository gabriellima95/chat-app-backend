package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Type     string `mapstructure:"DB_TYPE"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
}

func Load() {
	initializeViper(viper.GetViper())
}

// func bindEnv(v *viper.Viper, key ...string) {
// 	fmt.Println(key)

// 	err := v.BindEnv(key...)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

func initializeViper(env *viper.Viper) {
	// env.SetDefault("Env", "xxx")
	err := env.BindEnv("abc", "ABC", "XYZ")
	err = env.BindEnv("xyz", "XYZ")
	os.Setenv("ABC", "123")
	os.Setenv("XYZ", "321")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(env.Get("abc"))
	fmt.Println(env.Get("xyz"))

}
