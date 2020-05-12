package srcache

import "github.com/zhuyanxi/srcache/consistenthash"

const defaultLocalURL = "localhost:9099"
const defaultPathPrefix = "/_srcache/"
const defaultReplicate = 50
const defaultCapacity = 10000

// Server :
type Server struct {
	// localURL string

	// pathPrefix string

	// replicate int

	opts ServerOptions

	cache *SRCache

	peers *consistenthash.ConsistentHash

	hashFunc consistenthash.HashFunc

	// callback is the function that will query data from a outside data source
	// when the cache is missed
	callback callbackFunc
}

// ServerOptions :
type ServerOptions struct {
	LocalURL string

	CacheCapacity uint

	PathPrefix string

	Replicate int
}

type callbackFunc func(key string) ([]byte, error)
