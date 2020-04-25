package consistenthash

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestConsistentHash(t *testing.T) {
	hash256 := sha256.Sum256([]byte("123"))
	fmt.Println(hash256)
	hash256Str := hex.EncodeToString(hash256[:])
	fmt.Println(hash256Str)
	fmt.Printf("%x\n", hash256[:])
	hash256Decode, _ := hex.DecodeString(hash256Str)
	fmt.Println(hash256Decode)
	hashBig := binary.BigEndian.Uint32(hash256[:])
	fmt.Println(hashBig)
	hashSmall := binary.LittleEndian.Uint32(hash256[:])
	fmt.Println(hashSmall)

	// hash := NewConsistentHash(3, func(data []byte) int {
	// 	return int(crc32.ChecksumIEEE(data))
	// })
	hash := NewConsistentHash(4, func(data []byte) int {
		h1 := sha256.Sum256(data)
		// hash3 := crc32.ChecksumIEEE(data)
		hash32 := binary.BigEndian.Uint32(h1[:])
		return int(hash32)
	})

	nodes := []Node{
		{Addr: "2"},
		{Addr: "4"},
		{Addr: "6"},
	}

	hash.Add(nodes...)

	m := hash.ring
	for k, v := range m {
		logrus.Infof("k is %d, v is %s\n", k, v.Addr)
	}

	value1 := "10.192.168.10"
	n1 := hash.Get(value1)
	logrus.Infof("Value %s will be stored on node %s", value1, n1.Addr)
}
