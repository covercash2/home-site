package config

import (
	"errors"
	"github.com/covercash2/home-site/model"
	"github.com/spf13/viper"
)

// TODO add schema to config and person struct

var defaultPortNumber = "8081"

// Config is a structured representation
// of the config file
type Config struct {
	CSRFKey   [32]byte
	Me        model.Person
	Port      string
	StaticDir string
}

// ParseConfigFromFile takes a path to
// a file and returns the Config
// that represents the file contents
//TODO verify
// viper returns 0 values instead of errors
func ParseConfigFromFile(path string) (Config, error) {
	var config Config

	viper.SetDefault("port", defaultPortNumber)
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	me := model.Person{
		EmailAddress: viper.GetString("me.mail"),
		Name:         viper.GetString("me.name"),
		PhoneNumber:  viper.GetString("me.phone"),
	}

	key, err := validateKey(viper.GetString("CSRFKey"))
	if err != nil {
		return config, err
	}

	// TODO validate static directory

	config = Config{
		CSRFKey:   key,
		Me:        me,
		Port:      viper.GetString("port"),
		StaticDir: viper.GetString("staticDir"),
	}

	return config, nil
}

// validate csrf key
// takes a string and returns a [32]byte
// TODO shore up security
func validateKey(raw string) ([32]byte, error) {
	return validateKeyLength(raw)
}

func validateKeyLength(raw string) ([32]byte, error) {
	var retval [32]byte

	bytes := []byte(raw)
	if len(bytes) != 32 {
		return retval, errors.New(
			"given key is the wrong length",
		)
	}

	copy(retval[:], bytes)
	return retval, nil
}
