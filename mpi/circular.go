package mpi

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type CircularMPI struct {
	Dummy                     // the real data
	locQueues map[string]int  // The local queues
	pending   map[string]bool // Represent the pending writes
	hosts     []string        // A liste of "ip:port" of the nodes involved in the ring
	next      string          // "ip:port" of the next node in the ring
	me        string          // "ip:port" of the current node
	prev      string          // "ip:port" of the previous node in the ring
	logger    *log.Logger     // To give logging capabilities to the CircularMPI
	mutex     *sync.Mutex
}

func (m *CircularMPI) Register(s string) bool {
	_, ok := m.locQueues[s]
	if !ok {
		m.locQueues[s] = 1
		m.Dummy.Register(s)
		m.pending[s] = false
		return true
	}
	return false
}

func (m *CircularMPI) Flush() error {
	// we will flush the buffer if it's too long
	//m.mutex.Lock()
	//defer m.mutex.Unlock()
	to_send := make(map[string]Versionable, len(m.pending))
	m.prepareData(to_send)
	err := m.sendData(to_send)
	if err == nil {
		// We will flush the data
		m.pending = make(map[string]bool)
	}
	return err
}
func (m *CircularMPI) Write(s string, v Versionable) {
	m.Dummy.Write(s, v)
	m.pending[s] = true
}

func (m CircularMPI) sendData(to_send map[string]Versionable) error {
	if len(to_send) == 0 {
		return errors.New("No data to send")
	}
	var network bytes.Buffer
	enc := gob.NewEncoder(&network) // Will write to network.
	err := enc.Encode(to_send)
	if err != nil {
		return err
	}
	ur := strings.Join([]string{"http://", m.next}, "")
	_, err2 := http.PostForm(ur, url.Values{"data": {network.String()}})
	if err2 != nil {
		return err2
	}
	m.logger.Println("Sending Data to", m.next)
	return nil
}

func (m CircularMPI) prepareData(to_send map[string]Versionable) {
	// First we will clean the data to send from all the
	// local data
	for queue, _ := range m.locQueues {
		_, ok := to_send[queue]
		if ok {
			delete(to_send, queue)
		}
	}
	// we will populate the to_send variable
	// at every iteration we check if we haven't
	// already sent the previous variable
	for key, v := range m.pending {
		if v {
			_, ok := to_send[key]
			if !ok {
				to_send[key] = m.ReadFirst(key)
			}
		}
	}
}
func (m *CircularMPI) cleanLocalStreams(states map[string]Versionable) {
	for k := range m.locQueues {
		if _, ok := states[k]; ok {
			delete(states, k)
		}
	}
}
func (m *CircularMPI) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{'O', 'K'})
	go func() {
		states := make(map[string]Versionable)
		if r.FormValue("data") != "" {
			data := strings.NewReader(r.FormValue("data"))
			dec := gob.NewDecoder(data)
			err := dec.Decode(&states)
			if err != nil {
				m.logger.Println(err)
			}
		}
		m.cleanLocalStreams(states)
		m.prepareData(states)
		if len(states) == 0 {
			m.logger.Println("Nothing to send Going")
			return
		}
		// Ne need to clean local streams
		for key, v := range states {
			m.Dummy.Write(key, v)
		}
		m.logger.Println("receiving data from ", r.Host)
		m.sendData(states)
	}()
}

func NewCircularMPI(me string, hosts []string, logger *log.Logger) *CircularMPI {
	nxt := ""
	prv := ""
	for i, v := range hosts {
		if v == me {
			nxt = hosts[(i+1)%len(hosts)]
			j := (i - 1) % len(hosts)
			if j < 0 {
				j = j + len(hosts)
			}
			prv = hosts[j]
		}
	}
	return &CircularMPI{
		Dummy:     NewDummy(),
		locQueues: make(map[string]int, len(hosts)),
		pending:   make(map[string]bool),
		hosts:     hosts,
		next:      nxt,
		me:        me,
		prev:      prv,
		logger:    logger,
		mutex:     new(sync.Mutex)}
}
