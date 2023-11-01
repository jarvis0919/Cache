package singleflight

import (
	"sync"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}
type ReqBuff struct {
	mu sync.Mutex
	m  map[string]*call
}

func (r *ReqBuff) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	r.mu.Lock()
	if r.m == nil {
		r.m = make(map[string]*call)
	}
	if c, ok := r.m[key]; ok {
		r.mu.Unlock()
		c.wg.Wait()

		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	r.m[key] = c
	r.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	r.mu.Lock()
	delete(r.m, key)
	r.mu.Unlock()
	return c.val, c.err
}
