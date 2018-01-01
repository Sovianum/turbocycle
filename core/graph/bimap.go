package graph

import "fmt"

func newBiMap() *biMap {
	return &biMap{
		forward:  make(map[int]interface{}),
		backward: make(map[interface{}]int),
	}
}

type pair struct {
	key int
	val interface{}
}

type biMap struct {
	forward  map[int]interface{}
	backward map[interface{}]int
}

func (b *biMap) Length() int {
	return len(b.forward)
}

func (b *biMap) Add(key int, val interface{}) error {
	if _, ok := b.forward[key]; ok {
		return fmt.Errorf("key %d already presents", key)
	}
	if _, ok := b.backward[key]; ok {
		return fmt.Errorf("val %v already presents", val)
	}
	b.forward[key] = val
	b.backward[val] = key
	return nil
}

func (b *biMap) ContainsKey(key int) bool {
	_, ok := b.forward[key]
	return ok
}

func (b *biMap) ContainsVal(val interface{}) bool {
	_, ok := b.backward[val]
	return ok
}

func (b *biMap) Iterate() chan pair {
	var result = make(chan pair)
	var f = func() {
		for key, val := range b.forward {
			result <- pair{key: key, val: val}
		}
		close(result)
	}
	go f()

	return result
}

func (b *biMap) GetByKey(key int) (val interface{}, ok bool) {
	val, ok = b.forward[key]
	return
}

func (b *biMap) GetByVal(val interface{}) (key int, ok bool) {
	key, ok = b.backward[val]
	return
}

func (b *biMap) DeleteByVal(val interface{}) {
	var key, ok = b.backward[val]
	delete(b.backward, val)

	if ok {
		delete(b.forward, key)
	}
}

func (b *biMap) DeleteByKey(key int) {
	var val, ok = b.forward[key]
	delete(b.forward, key)

	if ok {
		delete(b.backward, val)
	}
}
