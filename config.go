package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

type Config struct {
	TurnSeconds int `yaml:"turn_seconds"`
	Port int `yaml:"port"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}
