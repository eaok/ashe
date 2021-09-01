package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eaok/ashe/config"
	"github.com/eaok/ashe/router"
	"github.com/lonelyevil/khl"
	"github.com/lonelyevil/khl/log_adapter/plog"
	"github.com/phuslu/log"
)

func init() {
	// Read config from configuration
	config.ReadConfig("./config/config.ini")
}

func main() {
	s := khl.New(config.Data.Token, plog.NewLogger(&log.Logger{
		Level: log.ErrorLevel,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}))

	// router.InitAction(s)
	router.Route(s)

	// need to open the socket
	err := s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the KHL session.
	s.Close()
}
