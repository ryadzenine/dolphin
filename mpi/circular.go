package mpi

import (
	"log"
	"sync"
)

type BaseCircularMPI struct {
	Dummy                     // the real data
	locQueues map[string]int  // The local queues
	pending   map[string]bool // Represent the pending writes
	hosts     []string        // A liste of "ip:port" of the nodes involved in the ring
	next      string          // "ip:port" of the next node in the ring
	me        string          // "ip:port" of the current node
	prev      string          // "ip:port" of the previous node in the ring
	logger    *log.Logger     // To give logging capabilities to the BaseCircularMPI
	mutex     *sync.Mutex
}

func (m *BaseCircularMPI) Register(s string) bool {
	_, ok := m.locQueues[s]
	if !ok {
		m.locQueues[s] = 1
		m.Dummy.Register(s)
		m.pending[s] = false
		return true
	}
	return false
}

func (m *BaseCircularMPI) Write(s string, v Versionable) {
	m.Dummy.Write(s, v)
	m.pending[s] = true
}

func (m BaseCircularMPI) prepareData(toSend map[string]Versionable) {
	// First we will clean the data to send from all the
	// local data
	if toSend == nil {
		toSend = make(map[string]Versionable)
	}
	for queue := range m.locQueues {
		_, ok := toSend[queue]
		if ok {
			delete(toSend, queue)
		}
	}
	// we will populate the toSend variable
	// at every iteration we check if we haven't
	// already sent the previous variable
	for key, v := range m.pending {
		if v {
			_, ok := toSend[key]
			if !ok {
				toSend[key] = m.ReadFirst(key)
			}
		}
	}
}

func (m *BaseCircularMPI) cleanLocalStreams(states map[string]Versionable) {
	for k := range m.locQueues {
		if _, ok := states[k]; ok {
			delete(states, k)
		}
	}
}
