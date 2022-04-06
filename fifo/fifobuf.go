package fifo

import (
	"sync"
)

type Buffer struct {
	mu   sync.Mutex
	set  map[interface{}]struct{} // for de-duplicating= entries.
	list []interface{} // list of all stored items in buffer (fifo)
	head chan interface{} // keeps head of the fifo list
}

func NewBuffer() (*Buffer, <-chan interface{}) {
	s := &Buffer{
		mu:   sync.Mutex{},
		set:  make(map[interface{}]struct{}),
		list: make([]interface{}, 0),
		head: make(chan  interface{}, 1),
	}

	return s, s.head
}

func (s *Buffer) Check() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.check()
}

func (s *Buffer) Head() <-chan interface{}{
	return s.head
}

func (s *Buffer) check() {
	if len(s.list) == 0 {
		return
	}
	select {
	case s.head <- s.list[0]:
		delete(s.set, s.list[0])
		if len(s.list) == 1 {
			s.list = make([]interface{}, 0)
		} else {
			s.list = s.list[1:]
		}
	default:
		// head is still full, do nothing.
	}
}

func (s *Buffer) Add(i interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.set[i]; exists {
		return
	}

	s.set[i] = struct{}{}
	s.list = append(s.list, i)
	s.check()
}
