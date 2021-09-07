package main

import (
	"fmt"
	"os"
	"time"

	"github.com/evan-buss/openbooks/mock"
)

func main() {
	ready := make(chan struct{})

	ircServer := mock.IrcServer{
		Port: ":6667",
	}
	go ircServer.Start(ready)
	<-ready

	greatGatsby, err := os.Open("great-gatsby.epub")
	if err != nil {
		panic(err)
	}

	dccConfig := mock.DccServer{
		Port:   ":6669",
		Reader: greatGatsby,
	}

	go dccConfig.Start(ready)
	<-ready

	searchResults, err := os.Open("SearchBot_results_for__the_great_gatsby.txt.zip")
	if err != nil {
		panic(err)
	}

	dccSearch := mock.DccServer{
		Port:   ":6668",
		Reader: searchResults,
	}

	go dccSearch.Start(ready)
	<-ready

	fmt.Println("waiting")
	time.Sleep(time.Hour * 24)
}
