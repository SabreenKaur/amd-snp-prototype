package main

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

const DefaultTTL = time.Minute

type Vcek_TTLCache struct {
	cache *ttlcache.Cache[string, []byte]
}

func NewVcek_TTLCache() *Vcek_TTLCache {
	return &Vcek_TTLCache{
		cache: ttlcache.New[string, []byte](),
	}
}
