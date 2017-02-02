package sift

import "container/list"

// not thread safe
type Lru interface {
	Add(interface{}, interface{})
	Get(interface{}) (interface{}, bool)
	RemoveOldest()
	SetOnEvicted(func(interface{}, interface{}))
}

type lruEntry struct {
	Value interface{}
	Key   interface{}
}

type ClassicLru struct {
	MaxLen    int
	ll        *list.List
	cache     map[interface{}]*list.Element
	OnEvicted func(key interface{}, value interface{}) //evicted callback used for free memory
}

func NewClassicLru(MaxLen int) Lru {
	return &ClassicLru{
		MaxLen:MaxLen,
		ll:list.New(),
		cache:make(map[interface{}]*list.Element),
	}
}

func (cl *ClassicLru)SetOnEvicted(f func(interface{}, interface{})) {
	cl.OnEvicted = f
}
func (cl *ClassicLru)Add(key, value interface{}) {
	if cl.cache == nil {
		cl.ll = list.New()
		cl.cache = make(map[interface{}]*list.Element)
	}
	if ee, ok := cl.cache[key]; ok {
		cl.ll.MoveToFront(ee)
		ee.Value.(*lruEntry).Value = value
		return
	}
	ee := cl.ll.PushFront(&lruEntry{value, key})
	cl.cache[key] = ee
	if cl.MaxLen != 0 && cl.MaxLen < cl.ll.Len() {
		cl.RemoveOldest()
	}
}

func (cl *ClassicLru)Get(key interface{}) (value interface{}, ok bool) {
	if cl.cache == nil {
		return
	}

	if ee, ok := cl.cache[key]; ok {
		cl.ll.MoveToFront(ee)
		return ee.Value.(*lruEntry).Value, true
	}
	return nil, false
}

func (cl *ClassicLru)RemoveOldest() {
	if cl.cache == nil {
		return
	}
	ee := cl.ll.Back()
	le := ee.Value.(*lruEntry)
	if ee != nil {
		cl.ll.Remove(ee)
		delete(cl.cache, le.Key)
		if cl.OnEvicted != nil {
			cl.OnEvicted(le.Key, le.Value)
		}
	}
}
