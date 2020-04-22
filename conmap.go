package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

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
	var k string
	var v []byte
	for {
		select {
		case it := <-c.chAdd:
			k = it.key
			v = it.value
			c.m[k] = &item{
				key:   k,
				value: v,
			}
			logrus.Println("add ", it.key+"---"+string(it.value), "to the map")
			// fmt.Println("add ", it.key+"---"+string(it.value), "to the map")
		case c.chGet <- c.m:
			fmt.Println("get ", k+"-"+string(v), "of the map")
		}
	}
}

func (c *ConMap) Add(key string, value []byte) {
	time.Sleep(time.Millisecond)
	it := &item{
		key:   key,
		value: value,
	}
	c.chAdd <- it
}

func (c *ConMap) AddSync(key string, value []byte) {
	time.Sleep(time.Millisecond)
	it := &item{
		key:   key,
		value: value,
	}
	c.m[key] = it
	fmt.Println("add ", it.key+"---"+string(it.value), "to the map")
}

func (c *ConMap) Get(key string) ([]byte, bool) {
	cm := <-c.chGet
	v, ok := cm[key]
	if ok {
		return v.value, ok
	}
	return nil, false
}

func (c *ConMap) Len() int {
	return len(c.m)
}
