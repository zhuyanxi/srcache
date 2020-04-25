package srcache

import (
	"github.com/sirupsen/logrus"
	"github.com/zhuyanxi/srcache/cache"
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
			addFlag := sr.info.lru.Add(k, v)
			sr.info.nbytes += int64(len(k)) + int64(len(v))
			if addFlag {
				logrus.Println("add ", ca.key+"---"+string(ca.value), "to the map")
			} else {
				logrus.Println("modify ", ca.key+"---"+string(ca.value), "in the map")
			}
		case sr.chGet <- sr.info.lru:
		}
	}
}

// Add :
func (sr *SRCache) Add(key string, value []byte) {
	it := &item{
		key:   key,
		value: value,
	}
	sr.chAdd <- it
}

// Get :
func (sr *SRCache) Get(key string) ([]byte, bool) {
	lru := <-sr.chGet
	v, ok := lru.Get(key)
	if ok {
		logrus.Println("Get ", key+"---"+string(v), "success.")
		return v, ok
	}
	logrus.Println("Get ", key, "fail.")
	return nil, false
}

// Len :
func (sr *SRCache) Len() int {
	return sr.info.lru.Len()
}

// Data :
func (sr *SRCache) Data() map[string][]byte {
	//lrum:=sr.info.lru.
	return nil
}
