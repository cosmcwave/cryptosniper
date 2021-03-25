package main

import (
	"context"
	"cryptosniper/backend"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var symbolsFlag = flag.String("s", "", "Specify the symbols to snipe")
var extFlag = flag.String("e", "", "Specify a list of extensions you would like to use")
var intervalFlag = flag.String("i", "1h", "Specify the timeframe (m/h/w/d/y)")

func main() {
	flag.Parse()

	cfg, err := backend.NewDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	svc := backend.New(cfg)
	svc.Out.Info().Msg("config loaded")
	svc.Out.Info().Msg("backend initialized")

	if *symbolsFlag == "" {
		svc.Error.Error().Msgf("there are no symbols specified to snipe")
		return
	}
	symbols := strings.Split(*symbolsFlag, ",")
	svc.AddSymbols(symbols...)

	extensions := strings.Split(*extFlag, ",")
	svc.AddExtensions(extensions...)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		select {
		case <-c:
			cancel()
		}
	}()
	svc.Start(ctx)
}
