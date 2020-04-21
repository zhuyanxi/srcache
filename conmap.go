package main

import "fmt"

type item struct {
	key   string
	value []byte
}

type ConMap struct {
	m     map[string]*item
	chGet chan map[string]*item
	chAdd chan *item
}

func NewConMap() *ConMap {
	cm := &ConMap{
		m:     make(map[string]*item),
		chGet: make(chan map[string]*item),
		chAdd: make(chan *item),
	}
	go cm.mux()
	return cm
}

func (c *ConMap) mux() {
	for {
		select {
		case it := <-c.chAdd:
			k := it.key
			v := it.value
			c.m[k] = &item{
				key:   k,
				value: v,
			}
			fmt.Println("add ", it, "to the map")
		case c.chGet <- c.m:
		}
	}
}

func (c *ConMap) Add(key string, value []byte) {
	it := &item{
		key:   key,
		value: value,
	}
	c.chAdd <- it
}

func (c *ConMap) Get(key string) ([]byte, bool) {
	cm := <-c.chGet
	v, ok := cm[key]
	if ok {
		return v.value, ok
	}
	return nil, false
}

func (c *ConMap) Test() {
	fmt.Println("test")
}
