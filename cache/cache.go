package cache

import "cryptosniper/cache/memorycache"

type Cache interface {
	Get(key string) interface{}
	Set(key string, value interface{})
}

type CacheType int

const (
	MemoryCache CacheType = iota
)

func New(cacheType CacheType) Cache {
	switch cacheType {
	case MemoryCache:
		return memorycache.New()
	}
	return nil
}