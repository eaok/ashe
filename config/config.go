package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ConfigData struct {
	// Debug      string `json:"Debug"`
	Token      string `json:"Token"`
	DebugToken string `json:"DebugToken"`
	Prefix     string `json:"Prefix"`
}

var Data = ConfigData{}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// read config.json and get a Token
func ReadConfig() error {
	fmt.Println("Reading from config file...")
	file, err := ioutil.ReadFile("./config/config.json")
	checkError(err)

	err = json.Unmarshal(file, &Data)
	checkError(err)

	// debug, err := base64.StdEncoding.DecodeString(Data.Debug)
	// checkError(err)
	// Data.Token = string(debug)

	// decode token
	token, err := base64.StdEncoding.DecodeString(Data.Token)
	checkError(err)
	Data.Token = string(token)
	debugToken, err := base64.StdEncoding.DecodeString(Data.DebugToken)
	checkError(err)
	Data.DebugToken = string(debugToken)
	Data.DebugToken = "1/MTA0NzY=/XLUOtvHDdNSlFZG7D9OJ+A=="

	// fmt.Println("Got Debug:  " + Data.Debug)
	fmt.Println("Got Token:  " + Data.Token)
	fmt.Println("Got DebugToken:  " + Data.DebugToken)
	fmt.Println("Got Prefix: " + Data.Prefix)

	return nil
}
