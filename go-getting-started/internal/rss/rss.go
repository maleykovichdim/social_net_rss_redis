package rss

import (
	"github.com/go-redis/redis/v8"
)

type Rss struct {
	Redis         *redis.Client
	IsInitialized bool
}

// New service implementation.
func New() *Rss {
	return &Rss{
		Redis:         nil,
		IsInitialized: false,
	}
}
