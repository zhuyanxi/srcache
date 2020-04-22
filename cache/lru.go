package cache

import (
	"container/list"

	"github.com/sirupsen/logrus"
)

// LRUCache a cache implement lru algorithm
type LRUCache struct {
	// capacity of the cache; 0 means unlimit capacity
	capacity uint

	// max memory of bytes occupy; 0 means unlimit
	// maxBytes int64

	// current memory occupy
	// nbytes int64

	// doubly linked list
	dll *list.List

	// map structure cache
	cache map[string]*list.Element
}

// item the data structure stored in doubly linked list
type item struct {
	key   string
	value []byte
}

// New init the cache
func New(cap uint) *LRUCache {
	return &LRUCache{
		capacity: cap,
		// maxBytes: maxBytes,
		dll:   list.New(),
		cache: make(map[string]*list.Element),
	}
}

// Add add an item to the cache
// return false means the key is already exist, the add method just modify the value
// 		  true means the key is not exist, the add method actually add the value to the cache
func (lc *LRUCache) Add(key string, val []byte) bool {
	// if cache hitted, move the hitted value to the front of doubly linked list
	// and change the cache value to the new value, then return
	if hit, ok := lc.cache[key]; ok {
		lc.dll.MoveToFront(hit)
		hit.Value.(*item).value = val
		return false
	}

	// if not hit, add the new item to the front of doubly linked list
	ele := lc.dll.PushFront(&item{
		key:   key,
		value: val,
	})
	lc.cache[key] = ele

	// delete the oldest value if current length is larger than maxcap
	if lc.capacity != 0 && lc.Len() > int(lc.capacity) {
		lc.removeOldest()
	}
	return true
}

// Get find the value of a given key
func (lc *LRUCache) Get(key string) (value []byte, ok bool) {
	if hit, o := lc.cache[key]; o {
		lc.dll.MoveToFront(hit)
		return hit.Value.(*item).value, o
	}
	return
}

func (lc *LRUCache) removeOldest() {
	if ele := lc.dll.Back(); ele != nil {
		k := ele.Value.(*item).key
		v := ele.Value.(*item).value
		lc.removeEle(ele)
		logrus.Println("Remove oldest element: ", k+"---"+string(v))
	}
}

func (lc *LRUCache) removeEle(ele *list.Element) {
	lc.dll.Remove(ele)
	delete(lc.cache, ele.Value.(*item).key)
}

// Len returns the number of items in the cache.
func (lc *LRUCache) Len() int {
	if lc.cache == nil {
		return 0
	}
	return lc.dll.Len()
}
