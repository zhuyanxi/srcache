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
	info  *cacheInfo
	chAdd chan *item
	chGet chan *cache.LRUCache
}

type item struct {
	key   string
	value []byte
}

type cacheInfo struct {
	capacity uint
	nbytes   int64
	lru      *cache.LRUCache
}

// NewSRCache init srcache instance
func NewSRCache(cap uint) *SRCache {
	cInfo := &cacheInfo{
		capacity: cap,
		lru:      cache.New(cap),
	}

	sr := &SRCache{
		info:  cInfo,
		chAdd: make(chan *item),
		chGet: make(chan *cache.LRUCache),
	}

	go sr.do()
	return sr
}

func (sr *SRCache) do() {
	var k string
	var v []byte
	for {
		select {
		case ca := <-sr.chAdd:
			k = ca.key
			v = ca.value
			sr.info.lru.Add(k, v)
			sr.info.nbytes += int64(len(k)) + int64(len(v))
		case sr.chGet <- sr.info.lru:
		}
	}
}

func (sr *SRCache) Add(key string, value []byte) {
	it := &item{
		key:   key,
		value: value,
	}
	sr.chAdd <- it
}

func (sr *SRCache) Get(key string) ([]byte, bool) {
	lru := <-sr.chGet
	v, ok := lru.Get(key)
	if ok {
		return v.([]byte), ok
	}
	return nil, false
}
