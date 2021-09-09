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

	IDChannelBS      string
	NameChannelBS    string
	IDChannelGroupBS string

	IDMsgRS string
	IDMsgBS string

	RoleRS11 int64
	RoleRS10 int64
	RoleRS9  int64
	RoleRS8  int64
	RoleRS7  int64
	RoleRS6  int64
	RoleRS5  int64
	RoleRS4  int64

	RoleBS8 int64
	RoleBS7 int64
	RoleBS6 int64
	RoleBS5 int64
	RoleBS4 int64
	RoleBS3 int64
	RoleBS2 int64
	RoleBS1 int64
}

var Data = ConfigData{}
var EmojiNum = [4]string{}
var BSRoleNum = [8]int64{}

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

	Data.IDChannelBS = conf.Section(appMod).Key("IDChannelBS").String()
	Data.NameChannelBS = conf.Section(appMod).Key("NameChannelBS").String()
	Data.IDChannelGroupBS = conf.Section(appMod).Key("IDChannelGroupBS").String()

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
	Data.IDMsgBS = conf.Section(appMod).Key("IDMsgBS").String()

	Data.RoleRS11 = conf.Section(appMod).Key("RoleRS11").MustInt64(343467)
	Data.RoleRS10 = conf.Section(appMod).Key("RoleRS10").MustInt64(343467)
	Data.RoleRS9 = conf.Section(appMod).Key("RoleRS9").MustInt64(343467)
	Data.RoleRS8 = conf.Section(appMod).Key("RoleRS8").MustInt64(343467)
	Data.RoleRS7 = conf.Section(appMod).Key("RoleRS7").MustInt64(343467)
	Data.RoleRS6 = conf.Section(appMod).Key("RoleRS6").MustInt64(343467)
	Data.RoleRS5 = conf.Section(appMod).Key("RoleRS5").MustInt64(343467)
	Data.RoleRS4 = conf.Section(appMod).Key("RoleRS4").MustInt64(343467)

	Data.RoleBS8 = conf.Section(appMod).Key("RoleBS8").MustInt64(343467)
	Data.RoleBS7 = conf.Section(appMod).Key("RoleBS7").MustInt64(343467)
	Data.RoleBS6 = conf.Section(appMod).Key("RoleBS6").MustInt64(343467)
	Data.RoleBS5 = conf.Section(appMod).Key("RoleBS5").MustInt64(343467)
	Data.RoleBS4 = conf.Section(appMod).Key("RoleBS4").MustInt64(343467)
	Data.RoleBS3 = conf.Section(appMod).Key("RoleBS3").MustInt64(343467)
	Data.RoleBS2 = conf.Section(appMod).Key("RoleBS2").MustInt64(343467)
	Data.RoleBS1 = conf.Section(appMod).Key("RoleBS1").MustInt64(343467)

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

	// test
	// Data.Token = "1/MTA0NzY=/XLUOtvHDdNSlFZG7D9OJ+A=="

	fmt.Println("Got Token:  " + Data.Token)
	fmt.Println("Got Prefix: " + Data.Prefix)
	fmt.Println("Got appMod:  " + Data.RunMod)

	// 给全局变量赋值
	EmojiNum[0] = Data.EmojiOne
	EmojiNum[1] = Data.EmojiTwo
	EmojiNum[2] = Data.EmojiThree
	EmojiNum[3] = Data.EmojiFour

	BSRoleNum[0] = Data.RoleBS1
	BSRoleNum[1] = Data.RoleBS2
	BSRoleNum[2] = Data.RoleBS3
	BSRoleNum[3] = Data.RoleBS4
	BSRoleNum[4] = Data.RoleBS5
	BSRoleNum[5] = Data.RoleBS6
	BSRoleNum[6] = Data.RoleBS7
	BSRoleNum[7] = Data.RoleBS8

	RSEmoji[Data.RoleRS4] = Data.EmojiFour
	RSEmoji[Data.RoleRS5] = Data.EmojiFive
	RSEmoji[Data.RoleRS6] = Data.EmojiSix
	RSEmoji[Data.RoleRS7] = Data.EmojiSeven
	RSEmoji[Data.RoleRS8] = Data.EmojiNeight
	RSEmoji[Data.RoleRS9] = Data.EmojiNine
	RSEmoji[Data.RoleRS10] = Data.EmojiTen
	RSEmoji[Data.RoleRS11] = Data.EmojiEleven

	RSEmoji[Data.RoleBS8] = Data.EmojiNeight
	RSEmoji[Data.RoleBS7] = Data.EmojiSeven
	RSEmoji[Data.RoleBS6] = Data.EmojiSix
	RSEmoji[Data.RoleBS5] = Data.EmojiFive
	RSEmoji[Data.RoleBS4] = Data.EmojiFour
	RSEmoji[Data.RoleBS3] = Data.EmojiThree
	RSEmoji[Data.RoleBS2] = Data.EmojiTwo
	RSEmoji[Data.RoleBS1] = Data.EmojiOne

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
