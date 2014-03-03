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
  *Dummy                   // the real data
  locQueues map[string]int // The local queues
  pending   []string       // Represent the pending writes
  hosts     []string       // A liste of "ip:port" of the nodes involved in the ring
  next      string         // "ip:port" of the next node in the ring
  me        string         // "ip:port" of the current node
  prev      string         // "ip:port" of the previous node in the ring
  logger    *log.Logger    // To give logging capabilities to the CircularMPI
  mutex     *sync.Mutex
}

func (m *CircularMPI) Register(s string) bool {
  _, ok := m.locQueues[s]
  if !ok {
    m.locQueues[s] = 1
    m.Dummy.Register(s)
    return true
  }
  return false
}
func (m *CircularMPI) Flush() {
  // we will flush the buffer if it's too long
  m.mutex.Lock()
  defer m.mutex.Unlock()
  to_send := make(map[string]Versionable, len(m.pending))
  m.prepareData(to_send)
  err := m.sendData(to_send)
  m.logger.Println(m.me, " flushing errors : ", err)
  if err == nil {
    // We will flush the data
    m.pending = make([]string, 0, len(m.locQueues))
  }
}
func (m *CircularMPI) Write(s string, v Versionable) {
  m.mutex.Lock()
  defer m.mutex.Unlock()
  m.Dummy.Write(s, v)
  m.pending = append(m.pending, s)
}

func (m CircularMPI) sendData(to_send map[string]Versionable) error {
  if len(to_send) == 0 {
    return errors.New("No data to send")
  }
  var network bytes.Buffer
  enc := gob.NewEncoder(&network) // Will write to network.
  //  gob.Register(Mock(1))
  err := enc.Encode(to_send)
  if err != nil {
    return err
  }
  ur := strings.Join([]string{"http://", m.next}, "")
  m.logger.Println(m.me, " Sending data ", to_send, " To : ", ur)
  _, err2 := http.PostForm(ur, url.Values{"data": {network.String()}})
  if err2 != nil {
    return err2
  }
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
  for _, v := range m.pending {
    _, ok := to_send[v]
    if !ok {
      to_send[v] = m.ReadFirst(v)
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
  data := strings.NewReader(r.FormValue("data"))
  dec := gob.NewDecoder(data)
  states := make(map[string]Versionable)
  err := dec.Decode(&states)
  m.logger.Println(m.me, " Receiving data", states, " from", r.RemoteAddr)
  if err != nil {
    m.logger.Println(err)
  }
  m.cleanLocalStreams(states)
  if len(states) == 0 {
    m.logger.Print(m.me, "  Nothing to send, done")
    return
  }
  // Ne need to clean local streams
  for key, v := range states {
    m.Dummy.Write(key, v)
  }
  go m.sendData(states)
  w.Write([]byte{'O', 'K'})
}

func NewCircularMPI(capacity int, me string, hosts []string, logger *log.Logger) *CircularMPI {
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
    Dummy:     NewDummy(capacity),
    locQueues: make(map[string]int, len(hosts)),
    pending:   make([]string, 0, 10),
    hosts:     hosts,
    next:      nxt,
    me:        me,
    prev:      prv,
    logger:    logger,
    mutex:     new(sync.Mutex)}
}
