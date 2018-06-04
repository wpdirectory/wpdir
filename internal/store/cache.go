package store

import (
	"fmt"
	"time"
)

func (s *Store) CacheSet(key string, object interface{}, expiry time.Duration) {
	s.cache.Set(key, object, expiry)
	fmt.Println(s.cache.Items())
}

func (s *Store) CacheGet(key string) (interface{}, bool) {
	fmt.Println(s.cache.Items())
	return s.cache.Get(key)
}

func (s *Store) CacheDelete(key string) {
	s.cache.Delete(key)
}
