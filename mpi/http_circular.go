package mpi

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type HTTPCircularMPI struct {
	BaseCircularMPI
}

func (m *HTTPCircularMPI) Flush() error {
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

func (m HTTPCircularMPI) sendData(to_send map[string]Versionable) error {
	if len(to_send) == 0 {
		return errors.New("No data to send")
	}
	data, err := encodeData(to_send)
	if err != nil {
		return err
	}
	ur := strings.Join([]string{"http://", m.next}, "")
	_, err2 := http.PostForm(ur, url.Values{"data": {data.String()}})
	if err2 != nil {
		return err2
	}
	m.logger.Println("Sending Data to", m.next)
	return nil
}

func (m *HTTPCircularMPI) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{'O', 'K'})
	go func() {
		states := make(map[string]Versionable)
		if r.FormValue("data") != "" {
			data := strings.NewReader(r.FormValue("data"))
			st, err := decodeData(data)
			if err != nil {
				m.logger.Println(err)
			}
			states = st
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

func NewHTTPCircularMPI(me string, hosts []string, logger *log.Logger) *HTTPCircularMPI {
	nxt, prv := getNxtPrev(me, hosts)
	return &HTTPCircularMPI{BaseCircularMPI{
		Dummy:     NewDummy(),
		locQueues: make(map[string]int, len(hosts)),
		pending:   make(map[string]bool),
		hosts:     hosts,
		next:      nxt,
		me:        me,
		prev:      prv,
		logger:    logger,
		mutex:     new(sync.Mutex)}}
}
