package utils

import (
	"log"
	"os"
	"encoding/json"
)

type Conf struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	AccessLog  string `json:"accesslog`
	ErrorLog   string `json:"errorlog`
	DBName     string `json:"dbname"`
	DBUser     string `json:"dbuser"`
	DBPassword string `json:"dbpassword"`
	DBHost     string `json:"dbhost"`
	DBProtocol string `json:"dbprotocol"`
	DBPort     int    `json:"dbport"`
	Log				 bool		`json:"log"`
}

var Config Conf

func init() {
	conf, err := getConfig("config.json")
	if err != nil {
		log.Panicln(err)
	}
	Config = conf
}

func getConfig(filename string) (Conf, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Conf{}, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Conf{}
	err = decoder.Decode(&config)
	if err != nil {
		return Conf{}, err
	}
	return config, nil
}
