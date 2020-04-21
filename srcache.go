package srcache

import (
	"github.com/zhuyanxi/srcache/cache"
)

// ISRCache define the srcache interface
// type ISRCache interface {
// 	Get() (value []byte, ok bool)
// 	Add(key string, value []byte)
// }

// SRCache wrap the lru cache
type SRCache struct {
	lru chan *cache.LRUCache
}

func NewSRCache() *SRCache {
	return nil
}
