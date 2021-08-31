package config

// debug
// const (
// 	IDChannelSelectRole = "3302070048274601"
// 	IDChannelRS11       = "9818481414209918"
// 	IDChannelRS10       = "5376021166530892"
// 	IDChannelRS9        = "9171279399459261"
// 	IDChannelRS8        = "2987960541506647"
// 	IDChannelRS7        = "2390191251391503"
// 	IDChannelRS6        = "1111"
// 	IDChannelRS5        = "2222"
// 	IDChannelRS4        = "3333"

// 	NameChannelSelectRole = "选择角色"
// 	NameChannelRS11       = "RS11"
// 	NameChannelRS10       = "RS10"
// 	NameChannelRS9        = "RS9"
// 	NameChannelRS8        = "RS8"
// 	NameChannelRS7        = "RS7"
// 	NameChannelRS6        = "RS6"
// 	NameChannelRS5        = "RS5"
// 	NameChannelRS4        = "RS4"

// 	IDChannelBL6 = "3679435737123478"
// 	IDChannelBL5 = "8691216906202449"
// 	IDChannelBL4 = "3908793693006626"

// 	IDMsgRS = "38ead55a-4087-41b1-adf2-47e3cc16f564"
// 	IDMxdBL = "902789b4-33ff-4555-8687-0ed7f93f6f29"

// 	RoleRS11 int64 = 377400
// 	RoleRS10 int64 = 377401
// 	RoleRS9  int64 = 377402
// 	RoleRS8  int64 = 377403
// 	RoleRS7  int64 = 377404
// 	RoleRS6  int64 = 377405
// 	RoleRS5  int64 = 377406
// 	RoleRS4  int64 = 377407

// 	RoleBL6 int64 = 377587
// 	RoleBL5 int64 = 377583
// 	RoleBL4 int64 = 377586
// )

//no debug
const (
	IDChannelSelectRole = "1037018848351596"
	IDChannelRS11       = "1677703796956707"
	IDChannelRS10       = "9468609038325390"
	IDChannelRS9        = "1216032623059191"
	IDChannelRS8        = "5677712030392850"
	IDChannelRS7        = "3724367550066488"
	IDChannelRS6        = "9872794175561463"
	IDChannelRS5        = "7677407365093177"
	IDChannelRS4        = "5617095043467552"

	NameChannelSelectRole = "选择角色"
	NameChannelRS11       = "RS11"
	NameChannelRS10       = "RS10"
	NameChannelRS9        = "RS9"
	NameChannelRS8        = "RS8"
	NameChannelRS7        = "RS7"
	NameChannelRS6        = "RS6"
	NameChannelRS5        = "RS5"
	NameChannelRS4        = "RS4"

	IDChannelBL6 = "3679435737123478"
	IDChannelBL5 = "8691216906202449"
	IDChannelBL4 = "3908793693006626"

	IDMsgRS = "f7857ce4-c5a8-408a-a3a5-96eac3df2718"
	IDMxdBL = "902789b4-33ff-4555-8687-0ed7f93f6f29"

	RoleRS11 int64 = 343469
	RoleRS10 int64 = 343468
	RoleRS9  int64 = 343467
	RoleRS8  int64 = 343463
	RoleRS7  int64 = 343462
	RoleRS6  int64 = 343461
	RoleRS5  int64 = 343460
	RoleRS4  int64 = 343458

	RoleBL6 int64 = 377587
	RoleBL5 int64 = 377583
	RoleBL4 int64 = 377586
)

// emoji constants
const (
	EmojiPointDown = "👇"
	EmojiRedCircle = "🔴"
	EmojiCheckMark = "✅"
	EmojiCrossMark = "❎"
	EmojiStopSign  = "🛑"
	EmojiOne       = "\u0031\uFE0F\u20E3"
	EmojiTwo       = "\u0032\uFE0F\u20E3"
	EmojiThree     = "\u0033\uFE0F\u20E3"
	EmojiFour      = "\u0034\uFE0F\u20E3"
	EmojiFive      = "\u0035\uFE0F\u20E3"
	EmojiSix       = "\u0036\uFE0F\u20E3"
	EmojiSeven     = "\u0037\uFE0F\u20E3"
	EmojiNeight    = "\u0038\uFE0F\u20E3"
	EmojiNine      = "\u0039\uFE0F\u20E3"
	EmojiTen       = "🔟"
	EmojiEleven    = "🚻"
)

var EmojiTest = EmojiOne + EmojiTwo + EmojiThree + EmojiFour + EmojiFive + EmojiSix + EmojiSeven + EmojiNeight + EmojiNine + EmojiTen +
	EmojiEleven + EmojiPointDown + EmojiRedCircle + EmojiCheckMark + EmojiCrossMark + EmojiStopSign

var EmojiNum = [4]string{
	EmojiOne,
	EmojiTwo,
	EmojiThree,
	EmojiFour,
}

var RSEmoji = map[int64]string{
	RoleRS4:  EmojiFour,
	RoleRS5:  EmojiFive,
	RoleRS6:  EmojiSix,
	RoleRS7:  EmojiSeven,
	RoleRS8:  EmojiNeight,
	RoleRS9:  EmojiNine,
	RoleRS10: EmojiTen,
	RoleRS11: EmojiEleven,
}
