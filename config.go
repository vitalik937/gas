package gas

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
}

func (c *Config) Load(appName string, configFileName string, recurse bool) error {
	viper.Reset()
	viper.SetConfigName(configFileName)

	if recurse {
		loadConfigFileRecursively()
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln("Fatal error reading "+appName+" config file:", err)
	}
	return err
}

func loadConfigFileRecursively() error {
	wd, err := os.Getwd()
	if err != nil {
		log.Panicln("An error occurred in gas while trying to get the calling app's working directory:", err)
	}

	arr := strings.Split(wd+"/", "/")
	i := len(arr) - 1
	for i > 1 {
		fp := strings.Join(arr[0:i], "/")
		log.Println("Attempting to load config file from:", fp)
		viper.AddConfigPath(fp)
		i--
	}
	return err
}

func GetString(path string) string {
	v := viper.GetString(path)
	if strings.HasPrefix(v, "ENV[") {
		return getEnv(v[4:len(v)-1], path)
	}
	return v
}

func GetInt(path string) int {
	v := viper.GetString(path)
	if strings.HasPrefix(v, "ENV[") {
		i, err := strconv.Atoi(getEnv(v[4:len(v)-1], path))
		if err != nil {
			log.Panicln("cannot cast env value to int for:", path)
		}
		return i
	}
	return viper.GetInt(path)
}

func GetBool(path string) bool {
	v := viper.GetString(path)
	if strings.HasPrefix(v, "ENV[") {
		b, err := strconv.ParseBool(getEnv(v[4:len(v)-1], path))
		if err != nil {
			log.Panicln("cannot cast env value to int for:", path)
		}
		return b
	}
	return viper.GetBool(path)
}

func getEnv(name string, path string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Panicln("env value cannot be blank:", path)
	}
	return v
}
