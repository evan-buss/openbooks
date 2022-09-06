package main

import (
	"fmt"
	"math/rand"
	"path"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/evan-buss/openbooks/server"
)

// Generate a random username to avoid IRC name collisions if multiple users are hosting
// at the same time.
func generateUserName() string {
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(int64(rand.Int()))
	return fmt.Sprintf("%s_%s", gofakeit.Adjective(), gofakeit.Noun())
}

// Update a server config struct from globalFlags
func bindGlobalServerFlags(config *server.Config) {
	config.Version = fmt.Sprintf("OpenBooks Server %s", ircVersion)
	config.UserName = globalFlags.UserName
	config.Log = globalFlags.Log
	config.Server = globalFlags.Server
	config.SearchBot = globalFlags.SearchBot
	config.EnableTLS = globalFlags.EnableTLS
}

// Make sure the server config has a valid rate limit.
func ensureValidRate(rateLimit int, config *server.Config) {

	// If user enters a limit that's too low, set to default of 10 seconds.
	if rateLimit < 10 {
		rateLimit = 10
	}

	config.SearchTimeout = time.Duration(rateLimit) * time.Second
}

func sanitizePath(basepath string) string {
	cleaned := path.Clean(basepath)
	if cleaned == "/" {
		return cleaned
	}
	return cleaned + "/"
}
