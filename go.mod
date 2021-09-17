module github.com/eaok/ashe

go 1.16

require (
	github.com/go-ini/ini v1.62.0
	github.com/lonelyevil/khl v0.0.8
	github.com/lonelyevil/khl/log_adapter/plog v0.0.8
	github.com/phuslu/log v1.0.74
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
)

replace github.com/lonelyevil/khl => ../khl

replace github.com/lonelyevil/khl/log_adapter/plog => ../khl/log_adapter/plog
