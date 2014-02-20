package mpi

import "container/list"
type DummyMessagesQueue struct {
    // represent the queues 
    queues map[string]*list.List
    // represent the capacity of the queues when each queue fill it's capacity
    // the older messages are removed 
    capacity int
}

func (m *DummyMessagesQueue) Queues() []string {
    names := make([]string, 0, len(m.queues))
    for key, _ := range m.queues {
        names = append(names , key)
    }
    return names 
}

func (m *DummyMessagesQueue) Register(name string) bool{
    _, ok := m.queues[name]
    if !ok {
        m.queues[name] = new(list.List)
        return true 
    }
    return false 
}

func (m *DummyMessagesQueue) Write(queue string, data Versionable){
    if _, ok:=m.queues[queue]; ok  {
        m.queues[queue].PushFront(data) 
    }else {
        ls := list.List{}
        ls.PushFront(data)
        if ls.Len() > m.capacity{
            ls.Remove(ls.Back())
        }
        m.queues[queue] = &ls
    }
}

func (m *DummyMessagesQueue) ReadFirst(queue string) (data Versionable){
    v, _ :=  m.queues[queue].Front().Value.(Versionable)
    return v  
}

func (m *DummyMessagesQueue) ReadFirstAll() (data map[string]Versionable){
    data = make(map[string]Versionable);
    for key, value := range m.queues {
        v, _ :=  value.Front().Value.(Versionable)
        data[key] = v;
    }
    return data
}

func (m *DummyMessagesQueue) ReadStates(versions map[string]int) map[string]Versionable {
    tmp := m.ReadFirstAll()
    data := make(map[string]Versionable)
    for key,v := range tmp {
        if v.Version() > versions[key]{
            data[key] = v
        }
    }
    return data 
}

func NewDummyMessagesQueue(capacity int) (queue DummyMessagesQueue) {
    queue = DummyMessagesQueue{make(map[string]*list.List), capacity}
    return 
}
