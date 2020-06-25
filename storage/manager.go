package storage

import (
	"sync"
)

var (
	DefaultManager = newManager()
)

type Manager interface {
	New(component string)
	Get(component string) Storage
}

type manager struct {
	sync.Map
}

func newManager() Manager {
	return &manager{}
}

func (m *manager) New(component string) {
	m.Store(component, newStorage())
}

func (m *manager) Get(component string) Storage {
	s, ok := m.Load(component)
	if ok {
		return s.(Storage)
	}
	return nil
}
