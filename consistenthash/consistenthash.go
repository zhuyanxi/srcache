package consistenthash

import (
	"encoding/json"
	"hash/crc32"
	"sort"
	"strconv"
)

var defaultHashFunc = crc32.ChecksumIEEE

// HashFunc :
type HashFunc func(data []byte) int

// ConsistentHash :
type ConsistentHash struct {
	// the hash function which computes the ring key
	hash HashFunc

	// virtual node numbers = number of node * replicate
	replicate int

	// the sorted keys is used for better performance
	sortedKeys []int

	// the hash ring
	ring map[int]Node
}

// NewConsistentHash :
func NewConsistentHash(replicate int, fn HashFunc) *ConsistentHash {
	ch := new(ConsistentHash)
	ch.replicate = replicate

	if fn == nil {
		ch.hash = func(data []byte) int {
			return int(defaultHashFunc(data))
		}
	} else {
		ch.hash = fn
	}

	ch.ring = make(map[int]Node)

	return ch
}

// Add : add nodes to the consistent hash ring
func (ch *ConsistentHash) Add(nodes ...Node) {
	for _, node := range nodes {
		for i := 0; i < ch.replicate; i++ {
			nodeJSON, _ := json.Marshal(node)
			tt := strconv.Itoa(i) + string(nodeJSON)
			hash := ch.hash([]byte(tt))
			ch.ring[hash] = node
			ch.sortedKeys = append(ch.sortedKeys, hash)
		}
	}
	sort.Ints(ch.sortedKeys)
}

// Get :
func (ch *ConsistentHash) Get(key string) Node {
	hash := int(ch.hash([]byte(key)))
	lenK := len(ch.sortedKeys)

	// return the index of first item that >= val, using binary search algorithm
	idx := sort.Search(lenK-1, func(i int) bool {
		return ch.sortedKeys[i] >= hash
	})

	return ch.ring[ch.sortedKeys[idx]]
}

// Node :
type Node struct {
	Addr string
}
