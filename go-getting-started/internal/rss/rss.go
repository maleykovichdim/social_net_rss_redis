package rss

import (
	"github.com/go-redis/redis/v8"
)

// Service contains the core business logic separated from the transport layer.
// You can use it to back a REST, gRPC or GraphQL API.
type Rss struct {
	Redis         *redis.Client
	IsInitialized bool
}

// New service implementation.
func New() *Rss {
	// cdc := branca.NewBranca(tokenKey)
	// cdc.SetTTL(uint32(tokenLifespan.Seconds()))

	return &Rss{
		Redis:         nil,
		IsInitialized: false,
	}
}
