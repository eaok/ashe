package config

import (
	"encoding/base64"
	"fmt"

	"github.com/go-ini/ini"
)

type ConfigData struct {
	Prefix string
	Token  string
	RunMod string

	EmojiPointDown string
	EmojiRedCircle string
	EmojiCheckMark string
	EmojiCrossMark string
	EmojiStopSign  string
	EmojiOne       string
	EmojiTwo       string
	EmojiThree     string
	EmojiFour      string
	EmojiFive      string
	EmojiSix       string
	EmojiSeven     string
	EmojiNeight    string
	EmojiNine      string
	EmojiTen       string
	EmojiEleven    string

	IDChannelSelectRole string
	IDChannelRS11       string
	IDChannelRS10       string
	IDChannelRS9        string
	IDChannelRS8        string
	IDChannelRS7        string
	IDChannelRS6        string
	IDChannelRS5        string
	IDChannelRS4        string

	NameChannelSelectRole string
	NameChannelRS11       string
	NameChannelRS10       string
	NameChannelRS9        string
	NameChannelRS8        string
	NameChannelRS7        string
	NameChannelRS6        string
	NameChannelRS5        string
	NameChannelRS4        string

	IDChannelBL6 string
	IDChannelBL5 string
	IDChannelBL4 string

	IDMsgRS string
	IDMxdBL string

	RoleRS11 int64
	RoleRS10 int64
	RoleRS9  int64
	RoleRS8  int64
	RoleRS7  int64
	RoleRS6  int64
	RoleRS5  int64
	RoleRS4  int64

	RoleBL6 int64
	RoleBL5 int64
	RoleBL4 int64
}

var Data = ConfigData{}
var EmojiNum = [4]string{}

//[roleID]emoji
var RSEmoji = map[int64]string{}

//[channelID]roleID
var ChanRole = map[string]int64{}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// read config.json and get a Token
func ReadConfig(path string) error {
	fmt.Println("Reading from config file...")
	conf, err := ini.Load(path)
	checkError(err)

	appMod := conf.Section("").Key("app_mode").String()
	Data.Token = conf.Section(appMod).Key("Token").String()
	Data.IDChannelSelectRole = conf.Section(appMod).Key("IDChannelSelectRole").String()
	Data.IDChannelRS11 = conf.Section(appMod).Key("IDChannelRS11").String()
	Data.IDChannelRS10 = conf.Section(appMod).Key("IDChannelRS10").String()
	Data.IDChannelRS9 = conf.Section(appMod).Key("IDChannelRS9").String()
	Data.IDChannelRS8 = conf.Section(appMod).Key("IDChannelRS8").String()
	Data.IDChannelRS7 = conf.Section(appMod).Key("IDChannelRS7").String()
	Data.IDChannelRS6 = conf.Section(appMod).Key("IDChannelRS6").String()
	Data.IDChannelRS5 = conf.Section(appMod).Key("IDChannelRS5").String()
	Data.IDChannelRS4 = conf.Section(appMod).Key("IDChannelRS4").String()
	Data.IDChannelBL6 = conf.Section(appMod).Key("IDChannelBL6").String()
	Data.IDChannelBL5 = conf.Section(appMod).Key("IDChannelBL5").String()
	Data.IDChannelBL4 = conf.Section(appMod).Key("IDChannelBL4").String()

	Data.NameChannelSelectRole = conf.Section(appMod).Key("NameChannelSelectRole").String()
	Data.NameChannelRS11 = conf.Section(appMod).Key("NameChannelRS11").String()
	Data.NameChannelRS10 = conf.Section(appMod).Key("NameChannelRS10").String()
	Data.NameChannelRS9 = conf.Section(appMod).Key("NameChannelRS9").String()
	Data.NameChannelRS8 = conf.Section(appMod).Key("NameChannelRS8").String()
	Data.NameChannelRS7 = conf.Section(appMod).Key("NameChannelRS7").String()
	Data.NameChannelRS6 = conf.Section(appMod).Key("NameChannelRS6").String()
	Data.NameChannelRS5 = conf.Section(appMod).Key("NameChannelRS5").String()
	Data.NameChannelRS4 = conf.Section(appMod).Key("NameChannelRS4").String()

	Data.IDMsgRS = conf.Section(appMod).Key("IDMsgRS").String()
	Data.IDMxdBL = conf.Section(appMod).Key("IDMxdBL").String()

	Data.RoleRS11 = conf.Section(appMod).Key("RoleRS11").MustInt64(343467)
	Data.RoleRS10 = conf.Section(appMod).Key("RoleRS10").MustInt64(343467)
	Data.RoleRS9 = conf.Section(appMod).Key("RoleRS9").MustInt64(343467)
	Data.RoleRS8 = conf.Section(appMod).Key("RoleRS8").MustInt64(343467)
	Data.RoleRS7 = conf.Section(appMod).Key("RoleRS7").MustInt64(343467)
	Data.RoleRS6 = conf.Section(appMod).Key("RoleRS6").MustInt64(343467)
	Data.RoleRS5 = conf.Section(appMod).Key("RoleRS5").MustInt64(343467)
	Data.RoleRS4 = conf.Section(appMod).Key("RoleRS4").MustInt64(343467)

	Data.RoleBL6 = conf.Section(appMod).Key("RoleBL6").MustInt64(343467)
	Data.RoleBL5 = conf.Section(appMod).Key("RoleBL5").MustInt64(343467)
	Data.RoleBL4 = conf.Section(appMod).Key("RoleBL4").MustInt64(343467)

	Data.Prefix = conf.Section("common").Key("prefix").String()
	Data.EmojiPointDown = conf.Section("common").Key("EmojiPointDown").String()
	Data.EmojiRedCircle = conf.Section("common").Key("EmojiRedCircle").String()
	Data.EmojiCheckMark = conf.Section("common").Key("EmojiCheckMark").String()
	Data.EmojiCrossMark = conf.Section("common").Key("EmojiCrossMark").String()
	Data.EmojiStopSign = conf.Section("common").Key("EmojiStopSign").String()
	Data.EmojiOne = conf.Section("common").Key("EmojiOne").String()
	Data.EmojiTwo = conf.Section("common").Key("EmojiTwo").String()
	Data.EmojiThree = conf.Section("common").Key("EmojiThree").String()
	Data.EmojiFour = conf.Section("common").Key("EmojiFour").String()
	Data.EmojiFive = conf.Section("common").Key("EmojiFive").String()
	Data.EmojiSix = conf.Section("common").Key("EmojiSix").String()
	Data.EmojiSeven = conf.Section("common").Key("EmojiSeven").String()
	Data.EmojiNeight = conf.Section("common").Key("EmojiNeight").String()
	Data.EmojiNine = conf.Section("common").Key("EmojiNine").String()
	Data.EmojiTen = conf.Section("common").Key("EmojiTen").String()
	Data.EmojiEleven = conf.Section("common").Key("EmojiEleven").String()

	token, err := base64.StdEncoding.DecodeString(Data.Token)
	checkError(err)
	Data.Token = string(token)
	Data.RunMod = appMod

	Data.Token = "1/MTA0NzY=/XLUOtvHDdNSlFZG7D9OJ+A=="

	fmt.Println("Got Token:  " + Data.Token)
	fmt.Println("Got Prefix: " + Data.Prefix)
	fmt.Println("Got appMod:  " + Data.RunMod)

	// 给全局变量赋值
	EmojiNum[0] = Data.EmojiOne
	EmojiNum[1] = Data.EmojiTwo
	EmojiNum[2] = Data.EmojiThree
	EmojiNum[3] = Data.EmojiFour

	RSEmoji[Data.RoleRS4] = Data.EmojiFour
	RSEmoji[Data.RoleRS5] = Data.EmojiFive
	RSEmoji[Data.RoleRS6] = Data.EmojiSix
	RSEmoji[Data.RoleRS7] = Data.EmojiSeven
	RSEmoji[Data.RoleRS8] = Data.EmojiNeight
	RSEmoji[Data.RoleRS9] = Data.EmojiNine
	RSEmoji[Data.RoleRS10] = Data.EmojiTen
	RSEmoji[Data.RoleRS11] = Data.EmojiEleven

	ChanRole[Data.IDChannelRS11] = Data.RoleRS11
	ChanRole[Data.IDChannelRS10] = Data.RoleRS10
	ChanRole[Data.IDChannelRS9] = Data.RoleRS9
	ChanRole[Data.IDChannelRS8] = Data.RoleRS8
	ChanRole[Data.IDChannelRS7] = Data.RoleRS7
	ChanRole[Data.IDChannelRS6] = Data.RoleRS6
	ChanRole[Data.IDChannelRS5] = Data.RoleRS5
	ChanRole[Data.IDChannelRS4] = Data.RoleRS4

	return nil
}
