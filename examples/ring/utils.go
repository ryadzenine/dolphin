package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Ring struct {
	Hosts []string
	Me    string
	Id    int
}

func (r *Ring) Sync() {
	if r.Id == 1 {
		WaitForTheSlaves(r)
	} else {
		WaitFormTheMaster(r)
	}
}

func NewRing(me int, cfFile string) (*Ring, error) {
	data, err := ioutil.ReadFile(cfFile)
	if err != nil {
		return nil, err
	}
	ring := new(Ring)
	ring.Id = me
	err = json.Unmarshal(data, ring)
	if err != nil {
		return nil, err
	}
	if me < 1 || me > len(ring.Hosts) {
		return nil, errors.New("Hostname out of the index")
	}
	ring.Me = ring.Hosts[me-1]
	return ring, nil
}

func logDest(base string, me int) (io.Writer, error) {
	return os.OpenFile(fmt.Sprint(base, "-", me, ".log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

func buildLogger(base string, me int) *log.Logger {
	f, err := logDest(base, me)
	if err != nil {
		panic(err)
	}
	return log.New(f, "Ring: ", log.LstdFlags)
}

func WaitForTheSlaves(ring *Ring) {
	var gr sync.WaitGroup
	gr.Add(len(ring.Hosts) - 1)
	ln, err := net.Listen("tcp", ring.Me)
	if err != nil {
		panic(err)
		return
	}
	defer ln.Close()
	j := 0
	for j < len(ring.Hosts)-1 {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Waiting for slaves: Acception Error: ", err)
		}
		j = j + 1
		go func(conn net.Conn) {
			gr.Done()
			fmt.Fprintf(conn, "ok")
			conn.Close()
		}(conn)
	}
	gr.Wait()
	time.Sleep(150 * time.Millisecond)
	for _, hst := range ring.Hosts {
		if hst != ring.Me {
			rep, _ := net.Dial("tcp", hst)
			fmt.Fprint(rep, "ok")
			rep.Close()
		}
	}
	//fmt.Println("All done Ready to Go")
}

func WaitFormTheMaster(ring *Ring) {
	rep, _ := net.Dial("tcp", ring.Hosts[0])
	for rep == nil {
		rep, _ = net.Dial("tcp", ring.Hosts[0])
	}
	fmt.Fprint(rep, "OK")
	rep.Close()
	var gr sync.WaitGroup
	gr.Add(1)
	ln, err := net.Listen("tcp", ring.Me)
	if err != nil {
		panic(err)
		return
	}
	defer ln.Close()
	conn, errr := ln.Accept()
	if errr != nil {
		fmt.Print(errr)
	}
	fmt.Fprint(conn, "OK")
	conn.Close()
}
