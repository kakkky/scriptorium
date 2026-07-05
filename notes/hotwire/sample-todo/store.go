package main

import (
	"sort"
	"sync"
)

type Todo struct {
	ID    int
	Title string
	Done  bool
}

type Store struct {
	mu     sync.RWMutex
	items  map[int]*Todo
	nextID int
}

func NewStore() *Store {
	return &Store{items: make(map[int]*Todo)}
}

func (s *Store) List() []*Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Todo, 0, len(s.items))
	for _, t := range s.items {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *Store) Get(id int) (*Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.items[id]
	return t, ok
}

func (s *Store) Create(title string) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	t := &Todo{ID: s.nextID, Title: title}
	s.items[s.nextID] = t
	return t
}

func (s *Store) Update(id int, title string, done bool) (*Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.items[id]
	if !ok {
		return nil, false
	}
	t.Title = title
	t.Done = done
	return t, true
}

func (s *Store) Toggle(id int) (*Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.items[id]
	if !ok {
		return nil, false
	}
	t.Done = !t.Done
	return t, true
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}

var store = NewStore()
