package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Uri      string
	Database string
}

func (c *Config) Read(configFile string) {
	if _, err := toml.DecodeFile(configFile, &c); err != nil {
		log.Fatal(err)
	}
}
