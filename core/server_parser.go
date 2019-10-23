package core

import (
	"sort"
	"strings"
	"time"
)

var prefixes = [...]string{
	"~",
	"&",
	"@",
	"%",
	"+",
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
	s.Time = time.Now()
	servers := strings.Split(data, " ")
	output := make([]string, 0)

	for _, user := range servers {
		for _, prefix := range prefixes {
			if strings.Contains(user, prefix) && strings.Index(user, prefix) == 0 && len(user) > 1 {
				output = append(output, user[1:])
			}
		}
	}
	sort.Strings(output)
	s.Servers = output
}

// GetServers returns the IRC book servers that are online
func GetServers(servers chan<- []string) {
	cacheIsOld := time.Now().Sub(serverCache.Time) > (time.Second * 30)
	if len(serverCache.Servers) == 0 || cacheIsOld {
		ircConn.GetUsers("ebooks")
		oldTime := serverCache.Time
		for serverCache.Time.Equal(oldTime) {
			time.Sleep(time.Millisecond * 500)
		}
	}
	servers <- serverCache.Servers
}
