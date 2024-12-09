package server

import "sync"

// implementation of the observer pattern for the servers

type ServerData struct {
	URL     string
	IsAlive bool
}

type Servers struct {
	Data      []ServerData
	observers []chan bool
	mu        sync.Mutex
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

func (s *Servers) NotifyObservers(message bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, observer := range s.observers {
		go func(ch chan bool) {
			ch <- message
		}(observer)
	}
}
