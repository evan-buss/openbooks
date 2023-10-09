package main

import (
	"context"
	"github.com/evan-buss/openbooks/mock"
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
	highwayOp := mock.NewOperator(&mock.Config{
		Server:  "localhost:6667",
		Name:    "search",
		Channel: "#ebooks",
	})

	undernetOp := mock.NewOperator(&mock.Config{
		Server:  "localhost:6667",
		Name:    "search2",
		Channel: "#bookz",
	})

	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		cancel()
	}()

	go highwayOp.StartListening(ctx)
	go undernetOp.StartListening(ctx)

	<-ctx.Done()
}
