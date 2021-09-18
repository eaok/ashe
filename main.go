package main

import (
	"net/http"
	_ "net/http/pprof"
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
	go func() {
		if config.Data.RunMod == "debug" {
			http.ListenAndServe(":8081", nil)
		} else {
			http.ListenAndServe(":8080", nil)
		}
	}()

	s := khl.New(config.Data.Token, plog.NewLogger(&log.Logger{
		Level:      log.WarnLevel,
		TimeFormat: "2006-01-02 15:04:05",
		Caller:     1,
		// Writer: &log.ConsoleWriter{
		// 	ColorOutput:    true,
		// 	QuoteString:    true,
		// 	EndWithMessage: true,
		// },
	}))

	router.InitAction(s)
	router.Route(s)

	// need to open the socket
	err := s.Open()
	if err != nil {
		s.Logger.Error().Err("", err).Msg("error opening connection")
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	s.Logger.Warn().Msg("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the KHL session.
	s.Close()
}
