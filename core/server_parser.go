package core

import (
	"sort"
	"strings"
)

var prefixes = map[byte]struct{}{
	'~': {},
	'&': {},
	'@': {},
	'%': {},
	'+': {},
}

type IrcServers struct {
	ElevatedUsers []string `json:"elevatedUsers"`
	RegularUsers  []string `json:"regularUsers"`
}

// ParseServers parses the complete list of IRC users to get the elevated users which in
// this case are the download servers
func ParseServers(rawString string) IrcServers {
	allServers := strings.Split(rawString, " ")

	servers := IrcServers{
		ElevatedUsers: make([]string, 0),
		RegularUsers:  make([]string, 0),
	}

	for _, name := range allServers {
		if len(name) > 1 {
			if _, exists := prefixes[name[0]]; exists {
				servers.ElevatedUsers = append(servers.ElevatedUsers, name[1:])
			} else {
				servers.RegularUsers = append(servers.RegularUsers, name)
			}
		}
	}

	sort.Slice(servers.ElevatedUsers, ignoreCaseSort(servers.ElevatedUsers))
	sort.Slice(servers.RegularUsers, ignoreCaseSort(servers.RegularUsers))

	return servers
}

func ignoreCaseSort(items []string) func(i, j int) bool {
	return func(i, j int) bool {
		return strings.ToLower(items[i]) < strings.ToLower(items[j])
	}
}
