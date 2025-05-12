package service

import (
	"sync"
	"sync/atomic"
)

type ServerPool struct {
	Servers []*Backend
	Index   uint32
	Mutex   sync.Mutex
}

func NewServerPool(servers []*Backend) *ServerPool {
	return &ServerPool{
		Servers: servers,
	}
}

func (s *ServerPool) NextServer() *Backend {
	count := uint32(len(s.Servers))
	for i := uint32(0); i < count; i++ {
		id := atomic.AddUint32(&s.Index, 1) % count
		server := s.Servers[id]
		if server.Alive {
			return server
		}
	}
	return nil
}

func (s *ServerPool) GetServer() []*Backend {
	return s.Servers
}
