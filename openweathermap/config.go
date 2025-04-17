package openweathermap

import (
	"weather/common"
)

type Config struct {
	Base  string `yaml:"baseURL"`
	City  string `yaml:"city"`
	Appid string `yaml:"appid"`
}

func NewConfig(path string) *Config {

	cfg := &Config{
		Base:  "https://api.openweathermap.org/data/2.5",
		City:  "London,UK",
		Appid: "00112233445566778899aabbccddeeff",
	}

	common.LoadConfig(path, cfg)

	return cfg
}
