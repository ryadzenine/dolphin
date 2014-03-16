package mpi

import (
  "log"
  "net/http"
  "os"
  "testing"
  "time"
)

func TestNewCircularMPI(t *testing.T) {
  me := "localhost:10525"
  hosts := []string{
    "localhost:10525",
    "localhost:10526",
    "localhost:10527"}
  logger := log.New(os.Stderr, "Job: ", log.Ldate)
  mp := NewCircularMPI(me, hosts, logger)
  if mp.next != "localhost:10526" || mp.prev != "localhost:10527" {
    t.Error("Wrong next node in the ring")
  }
  me = "localhost:10526"
  hosts = []string{
    "localhost:10525",
    "localhost:10526",
    "localhost:10527"}
  logger = log.New(os.Stderr, "Job: ", log.Ldate)
  mp = NewCircularMPI(me, hosts, logger)
  if mp.next != "localhost:10527" || mp.prev != "localhost:10525" {
    t.Error("Wrong next node in the ring")
  }
  me = "localhost:10527"
  hosts = []string{
    "localhost:10525",
    "localhost:10526",
    "localhost:10527"}
  logger = log.New(os.Stderr, "Job: ", log.Ldate)
  mp = NewCircularMPI(me, hosts, logger)
  if mp.next != "localhost:10525" || mp.prev != "localhost:10526" {
    t.Error("Wrong next node in the ring")
  }
}
func buildMP(i int) *CircularMPI {
  hosts := []string{
    "127.100.1.1:8080",
    "127.100.1.2:8081",
    "127.100.1.3:8082"}
  logger := log.New(os.Stderr, "Job: ", log.Ldate)
  return NewCircularMPI(hosts[i], hosts, logger)
}
func TestRegister(t *testing.T) {
  mp := buildMP(1)
  mp.Register("stream1")
  v, ok := mp.locQueues["stream1"]
  if v != 1 || !ok {
    t.Error("Registering failled")
  }
}

func TestWrite(t *testing.T) {
  mp1 := buildMP(1)
  mp1.Register("stream1")
  mp1.Write("stream1", Mock(1))
  if _, ok := mp1.Dummy["stream1"]; !ok || len(mp1.pending) != 1 {
    t.Error("Write failled")
  }
}

func TestPrepareData(t *testing.T) {
  mp1 := buildMP(1)
  mp1.Register("stream1")
  mp1.Write("stream1", Mock(3))
  to_send := map[string]Versionable{
    "stream":  Mock(2),
    "stream1": Mock(1)}
  mp1.prepareData(to_send)
  if v, ok := to_send["stream"]; !ok || v.Version() != 2 {
    t.Error("Function prepare data, deletes data the it should not, check it")
  }
  if v, ok := to_send["stream1"]; !ok || v.Version() != 3 {
    t.Error("Function prepareData, doesn't update data with pending as it's supposed to.")
  }
}

func TestSendReceiveData(t *testing.T) {
  mp1 := buildMP(0)
  mp1.Register("stream1")
  mp2 := buildMP(1)
  mp2.Register("stream2")
  mp3 := buildMP(2)
  mp3.Register("stream0")
  s1 := http.NewServeMux()
  s1.HandleFunc("/", mp1.MessagesHandler)
  go http.ListenAndServe(mp1.me, s1)
  s2 := http.NewServeMux()
  s2.HandleFunc("/", mp2.MessagesHandler)
  go http.ListenAndServe(mp2.me, s2)
  s3 := http.NewServeMux()
  s3.HandleFunc("/", mp3.MessagesHandler)
  go http.ListenAndServe(mp3.me, s3)
  mp1.Write("stream1", Mock(1))
  mp1.Register("stream4")
  mp1.Write("stream4", Mock(4))
  mp1.Flush()
  <-time.After(1500 * time.Microsecond)
  if len(mp1.pending) != 0 {
    t.Error("pending should be cleared after a send")
  }

}
