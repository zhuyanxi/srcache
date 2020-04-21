package core

import (
	"github.com/zhuyanxi/sr-cache/cache"
)

// ISRCache define the srcache interface
// type ISRCache interface {
// 	Get() (value []byte, ok bool)
// 	Add(key string, value []byte)
// }

// SRCache wrap the lru cache
type SRCache struct {
	chInfo chan *cacheInfo
	chLRU  chan *cache.LRUCache
}

type cacheInfo struct {
	capacity uint
	nbytes   int64
}

// NewSRCache init srcache instance
func NewSRCache() *SRCache {
	return &SRCache{
		chInfo: make(chan *cacheInfo),
		chLRU:  make(chan *cache.LRUCache),
	}
}

func (sr *SRCache) Add(key string, value []byte) {
	// if sr.lru == nil {
	// 	sr.lru = cache.New(10000)
	// }
}

func (sr *SRCache) Get(key string) ([]byte, bool) {
	var result []byte
	var ok bool
	select {
	case <-sr.chInfo:

	}
	return nil, false
}
