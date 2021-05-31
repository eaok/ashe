package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type configStruct struct {
	Token  string `json:"Token"`
	Prefix string `json:"Prefix"`
}

var (
	Token  string
	Prefix string
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// read config.json and get a Token
func ReadConfig() string {
	fmt.Println("Reading from config file...")
	file, err := ioutil.ReadFile("./config/config.json")
	checkError(err)

	config := configStruct{}
	err = json.Unmarshal(file, &config)
	checkError(err)

	Token = config.Token
	Prefix = config.Prefix
	fmt.Println("Got Token:  " + config.Token)
	fmt.Println("Got Prefix: " + config.Prefix)

	return config.Token
}
