package server

import "github.com/evan-buss/openbooks/core"

type Repository struct {
	servers core.IrcServers
}

func NewRepository() *Repository {
	return &Repository{servers: core.IrcServers{}}
}
