package main

import (
	"context"
	"github.com/evan-buss/openbooks/mock"
	"log"
	"os"
	"os/signal"
)

func main() {
	// Start IRC server container
	ctx, cancel := context.WithCancel(context.Background())
	//container, err := mock.StartIrcServer(ctx)
	//if err != nil {
	//	panic(err)
	//}

	//Start mock book operator
	operator := mock.NewOperator(&mock.Config{
		Server: "localhost:6667",
	})

	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		cancel()
	}()

	log.Fatal(operator.StartListening(ctx))
}
