package types

import "sync"

type ConcurrentSlice struct {
	sync.RWMutex
	ipExistMap map[string]bool
	slice      []string
}

func (c *ConcurrentSlice) Add(ips ...string) {
	c.Lock()
	defer c.Unlock()
	for _, each := range ips {
		if !c.ipExistMap[each] {
			c.slice = append(c.slice, each)
			c.ipExistMap[each] = true
		}
	}
}

func (c *ConcurrentSlice) GetSlice() *[]string {
	c.RLock()
	defer c.RUnlock()
	result := make([]string, len(c.slice))
	copy(result, c.slice)
	return &result
}

func NewConcurrentSlice() *ConcurrentSlice {
	return &ConcurrentSlice{
		ipExistMap: make(map[string]bool),
		slice:      make([]string, 0),
	}
}

type Storage struct {
	sync.RWMutex
	storage map[string]*ConcurrentSlice
}

func (s *Storage) Get(userId string) *ConcurrentSlice {
	s.RLock()
	defer s.RUnlock()
	return s.storage[userId]
}

func (s *Storage) Set(userId string, ips []string) {
	s.Lock()
	defer s.Unlock()
	conSlice, ok := s.storage[userId]
	if !ok {
		slice := NewConcurrentSlice()
		slice.Add(ips...)
		s.storage[userId] = slice
		return
	}

	conSlice.Add(ips...)
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]*ConcurrentSlice),
	}
}
