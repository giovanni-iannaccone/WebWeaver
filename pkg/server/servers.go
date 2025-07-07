package server

import "sync"

// implementation of the observer pattern for the servers

type Servers struct {
	Active      []string
	Inactive	[]string
	observers 	[]chan bool
	mu        	sync.Mutex
}

func (s *Servers) AddObserver(obs chan bool) {
	s.mu.Lock()
	s.observers = append(s.observers, obs)
	s.mu.Unlock()
}

func (s *Servers) RemoveObserver(obs chan bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, observer := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			return
		}
	}
}

func (s *Servers) NotifyObservers() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, observer := range s.observers {
		go func(ch chan bool) {
			ch <- true
		}(observer)
	}
}
