package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBName     string `json:"db_name"`
}

var Conf Config

func LoadConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &Conf)
	if err != nil {
		return err
	}
	return nil
}
