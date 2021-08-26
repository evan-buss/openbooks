package core

import (
	"sort"
	"strings"
	"time"
)

var prefixes = map[byte]struct{}{
	'~': {},
	'&': {},
	'@': {},
	'%': {},
	'+': {},
}

// ServerCache maintains a list of download servers that are available
// Time ensures that the cache is never too far out of date
type ServerCache struct {
	Servers []string
	Time    time.Time
}

// ParseServers parses the complete list of IRC users to get the elevates users which in
// this case are the download servers
func (s *ServerCache) ParseServers(data string) {
	servers := strings.Split(data, " ")
	output := make([]string, 0)

	for _, name := range servers {
		if _, ok := prefixes[name[0]]; ok && len(name) > 1 {
			output = append(output, name[1:])
		}
	}

	sort.Slice(output, func(i, j int) bool {
		return strings.ToLower(output[i]) < strings.ToLower(output[j])
	})

	s.Servers = output
	s.Time = time.Now()
}

// GetServers returns the IRC book servers that are online
// TODO: Look into a cleaner way of doing this
func GetServers(servers chan<- []string) {
	cacheIsOld := time.Since(serverCache.Time) > (time.Minute * 2)
	if len(serverCache.Servers) == 0 || cacheIsOld {
		ircConn.GetUsers("ebooks")
		oldTime := serverCache.Time
		for serverCache.Time.Equal(oldTime) {
			time.Sleep(time.Millisecond * 500)
		}
	}
	servers <- serverCache.Servers
}
