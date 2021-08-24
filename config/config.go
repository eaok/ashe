package config

import (
	"encoding/base64"
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

	// decode token
	token, err := base64.StdEncoding.DecodeString(config.Token)
	checkError(err)

	Token = string(token)
	Prefix = config.Prefix
	fmt.Println("Got Token:  " + Token)
	fmt.Println("Got Prefix: " + Prefix)

	return config.Token
}
