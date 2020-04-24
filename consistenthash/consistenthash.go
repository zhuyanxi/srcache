package consistenthash

import (
	"encoding/json"
	"hash/crc32"
	"sort"
	"strconv"
)

type hashFn func(data []byte) uint32

// ConsistentHash :
type ConsistentHash struct {
	hash hashFn

	replicate int

	// the sorted keys is used for better performance
	sortedKeys []int

	// the hash ring
	ring map[int]Node
}

// NewConsistentHash :
func NewConsistentHash(replicate int, fn hashFn) *ConsistentHash {
	ch := new(ConsistentHash)
	ch.replicate = replicate
	if fn == nil {
		ch.hash = crc32.ChecksumIEEE
	} else {
		ch.hash = fn
	}

	return ch
}

// Add : add nodes to the consistent hash ring
func (ch *ConsistentHash) Add(nodes ...Node) {
	for _, node := range nodes {
		for i := 0; i < ch.replicate; i++ {
			nodeJSON, _ := json.Marshal(node)
			hash := int(ch.hash([]byte(string(nodeJSON) + strconv.Itoa(i))))
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
	IP string
}
