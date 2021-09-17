package cli

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/irc"
	"github.com/evan-buss/openbooks/util"
)

func registerShutdown(conn *irc.Conn, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		conn.Disconnect()
		cancel()
		os.Exit(0)
	}()
}

func instantiate(config Config) *irc.Conn {
	fmt.Printf("Connecting to %s... ", config.Server)
	conn := irc.New(config.UserName, "OpenBooks CLI")
	core.Join(conn, config.Server)
	fmt.Println("✅")
	return conn
}

func (config *Config) setupLogger(handler core.EventHandler) io.Closer {
	logger, file, err := util.CreateLogFile(config.UserName, config.Dir)
	if err != nil {
		log.Fatalf("Error setting up logger: %s\n", err)
	}
	handler[core.Message] = func(text string) {
		logger.Println(text)
	}

	return file
}
