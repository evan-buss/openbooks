package main

import "github.com/samber/lo"

type Store interface {
	Authenticate(username, password string) bool
}

type ConfigStore struct {
	config *Config
}

func NewConfigStore(config *Config) *ConfigStore {
	return &ConfigStore{config: config}
}

func (c *ConfigStore) Authenticate(username, password string) bool {
	_, ok := lo.Find[User](c.config.Users, func(user User) bool {
		return user.Username == username && user.Password == password
	})

	return ok
}
