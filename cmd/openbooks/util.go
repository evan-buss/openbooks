package main

import (
	"path"
	"time"

	"github.com/evan-buss/openbooks/server"
)

// Update a server config struct from globalFlags
func bindGlobalServerFlags(config *server.Config, globalFlags GlobalFlags) {
	config.UserAgent = globalFlags.UserAgent
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
