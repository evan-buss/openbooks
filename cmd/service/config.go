package main

type Config struct {
	Debug    bool        `json:"debug,omitempty"`
	Port     string      `json:"port,omitempty"`
	Users    []User      `json:"users,omitempty"`
	Indexing IndexConfig `json:"indexing"`
	Irc      IrcConfig   `json:"irc"`
	Key      string      `json:"key"`
}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type IndexConfig struct {
	SourceDirectory string `json:"source_directory,omitempty"`
	StoreDirectory  string `json:"store_directory,omitempty"`
}

type IrcConfig struct {
	Url           string   `json:"url,omitempty"`
	SearchCommand string   `json:"search_command,omitempty"`
	Servers       []string `json:"servers,omitempty"`
}
