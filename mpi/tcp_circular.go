package mpi

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"
)

type CircularMPI struct {
	BaseCircularMPI
	connType string
}

func (m CircularMPI) sendData(toSend map[string]Versioner) error {
	if len(toSend) == 0 {
		return errors.New("no data to send")
	}
	data, err := encodeData(toSend)
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

func (m *CircularMPI) Flush() error {
	// we will flush the buffer if it's too long
	//m.mutex.Lock()
	//defer m.mutex.Unlock()
	toSend := make(map[string]Versioner, len(m.pending))
	m.prepareData(toSend)
	err := m.sendData(toSend)
	if err == nil {
		// We will flush the data
		m.pending = make(map[string]bool)
	}
	return err
}

func (m *CircularMPI) handleMessage(conn io.ReadWriteCloser) {
	defer conn.Close()
	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		m.logger.Println(err)
		conn.Write([]byte{'N', 'O'})
		return
	}
	conn.Write([]byte{'O', 'K'})

	go func() {
		states := make(map[string]Versioner)
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
		m.sendData(states)
	}()
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
		go m.handleMessage(conn)
	}
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
