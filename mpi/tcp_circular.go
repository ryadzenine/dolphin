package mpi

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"sync"
)

type CircularMPI struct {
	BaseCircularMPI
	connType string
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

func (m CircularMPI) sendData(to_send map[string]Versionable) error {
	if len(to_send) == 0 {
		return errors.New("No data to send")
	}
	data, err := encodeData(to_send)
	if err != nil {
		return err
	}
	conn, err := net.Dial(m.connType, m.next)
	if err != nil {
		m.logger.Println(err)
		return err
	}
	defer conn.Close()
	_, err2 := conn.Write(data.Bytes())
	if err2 != nil {
		m.logger.Println(err2)
		return err2
	}
	m.logger.Println("Sending Data to", m.next)
	return nil
}
func (m *CircularMPI) ListenAndServe() {
	ln, err := net.Listen(m.connType, m.me)
	if err != nil {
		m.logger.Fatal("Error listening:", err.Error())
		return
	}
	// Close the listener when the application closes.
	defer ln.Close()
	for {
		// Listen for an incoming connection.
		conn, err := ln.Accept()

		if err != nil {
			m.logger.Println("Error accepting: ", err.Error())
		}
		// Handle connections in a new goroutine.
		go m.MessagesHandler(conn)
	}
}
func (m *CircularMPI) MessagesHandler(conn net.Conn) {
	defer conn.Close()
	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		m.logger.Println(err)
		conn.Write([]byte{'N', 'O'})
		return
	}
	conn.Write([]byte{'O', 'K'})

	go func() {
		states := make(map[string]Versionable)
		data := bytes.NewBuffer(buf)
		st, err := decodeData(data)
		if err != nil {
			m.logger.Println(err)
		}
		states = st
		m.logger.Println(m.me, " Received ")
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
		m.logger.Println("receiving data from ", conn.RemoteAddr())
		m.sendData(states)
	}()
}

func NewCircularMPI(conn string, me string, hosts []string, logger *log.Logger) *CircularMPI {
	nxt, prv := getNxtPrev(me, hosts)
	return &CircularMPI{BaseCircularMPI{
		Dummy:     NewDummy(),
		locQueues: make(map[string]int, len(hosts)),
		pending:   make(map[string]bool),
		hosts:     hosts,
		next:      nxt,
		me:        me,
		prev:      prv,
		logger:    logger,
		mutex:     new(sync.Mutex)}, conn}
}
