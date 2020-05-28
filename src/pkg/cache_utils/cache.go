package cache_utils

import (
	time "time"

	axon_utils "github.com/jeroenvanmaanen/dendrite/src/pkg/axon_utils"
	go_cache "github.com/patrickmn/go-cache"
)

type GoCache interface {
	axon_utils.Cache
}

var cache GoCache

type goCache struct {
	cache *go_cache.Cache
}

func (g goCache) Get(s string) (interface{}, bool) {
	return g.cache.Get(s)
}

func (g goCache) Put(s string, i interface{}) {
	g.cache.SetDefault(s, i)
}

func (g goCache) Delete(s string) {
	g.cache.Delete(s)
}

func InitializeCache() {
	expireDuration, _ := time.ParseDuration("5m")
	cleanupInterval, _ := time.ParseDuration("10s")
	cache = goCache{
		cache: go_cache.New(expireDuration, cleanupInterval),
	}
	_ = axon_utils.SetCache(cache)
}
