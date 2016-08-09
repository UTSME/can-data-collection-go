package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Config struct {
	SerialPort string
	BaudRate   uint
}

var defaultConfig Config = Config{
	SerialPort: "COM1",
	BaudRate:   115200,
}

func (c *Config) LoadFromFile(filepath string) error {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	var fileConfig map[string]interface{}

	err = json.Unmarshal(b, &fileConfig)
	if err != nil {
		return errors.New("Error parsing config (JSON): " + err.Error())
	}

	for k, v := range fileConfig {
		switch k {
		case "SerialPort":
			c.SerialPort = v.(string)
		case "BaudRate":
			c.BaudRate = uint(v.(float64))
		default:
			return errors.New("Unknown config: " + k)
		}
	}

	//We made it fam.
	return nil
}
