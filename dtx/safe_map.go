package dtx

import (
	"iter"
	"maps"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]V
}

func NewSafeMap[K comparable, V any](cap int) *SafeMap[K, V] {
	return &SafeMap[K, V]{
		mu: sync.Mutex{},
		m:  make(map[K]V, cap),
	}
}

func (sm *SafeMap[K, V]) Clone() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		mu: sync.Mutex{},
		m:  maps.Clone(sm.m),
	}
}

func (sm *SafeMap[K, V]) ToMap() map[K]V {
	return maps.Clone(sm.m)
}

func (sm *SafeMap[K, V]) Len() int {
	return len(sm.m)
}

func (sm *SafeMap[K, V]) Set(k K, v V) {
	sm.mu.Lock()
	sm.m[k] = v
	sm.mu.Unlock()
}
func (sm *SafeMap[K, V]) Delete(k K) {
	sm.mu.Lock()
	delete(sm.m, k)
	sm.mu.Unlock()
}
func (sm *SafeMap[K, V]) Get(k K) (v V, ok bool) {
	sm.mu.Lock()
	v, ok = sm.m[k]
	sm.mu.Unlock()
	return v, ok
}

func (sm *SafeMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		sm.mu.Lock()

		// snapshot
		tmp := make(map[K]V, len(sm.m))
		for k, v := range sm.m {
			tmp[k] = v
		}

		sm.mu.Unlock()

		for _, v := range tmp {
			if !yield(v) {
				return
			}
		}
	}
}

func (sm *SafeMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		sm.mu.Lock()

		// snapshot
		tmp := make(map[K]V, len(sm.m))
		for k, v := range sm.m {
			tmp[k] = v
		}

		sm.mu.Unlock()

		for k, v := range tmp {
			if !yield(k, v) {
				return
			}
		}
	}
}
