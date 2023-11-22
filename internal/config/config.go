package config

import (
	"fmt"
	"log"

	"github.com/pelletier/go-toml"
)

var (
	globalConfig GlobalConfig
)

// Global Configuration
type GlobalConfig struct {
	Database Database `toml:"database"`
	Server   Server   `toml:"server"`
}

// DB configuration
type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	DBname   string `toml:"dbname"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

// server configuration
type Server struct {
	Address      string `toml:"address"`
	ReadTimeOut  int    `toml:"read_time_out"`
	WriteTimeOut int    `toml:"write_time_out"`
}

// Setter method for GlobalConfig
func SetConfig(cfg GlobalConfig) {
	globalConfig = cfg
}

// Getter method for GlobalConfig
func GetConfig() GlobalConfig {
	return globalConfig
}

// Loading the values from default.toml and assigning them as part of GlobalConfig struct
func InitGlobalConfig() error {
	config, err := toml.LoadFile("./../config/defaults.toml")
	fmt.Println("Err : ", err)
	if err != nil {
		log.Printf("Error while loading defaults.toml file : %v ", err)
		return err
	}

	var appConfig GlobalConfig
	err = config.Unmarshal(&appConfig)
	if err != nil {
		log.Printf("Error while unmarshalling config : %v", err)
		return err
	}

	SetConfig(appConfig)
	return nil
}
