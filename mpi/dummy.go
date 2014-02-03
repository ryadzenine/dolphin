package mpi

import "container/list"

type DummyMessagesQueue struct {
    // represent the queues 
    queues map[string]*list.List
    // represent the capacity of the queues when each queue fill it's capacity
    // the older messages are removed 
    capacity int
}

func (m *DummyMessagesQueue) Write(queue string, data interface{}){
    if _, ok:=m.queues[queue]; ok  {
        m.queues[queue].PushFront(data) 
    }else {
        ls := list.List{}
        ls.PushFront(data)
        m.queues[queue] = &ls
    }
}

func (m *DummyMessagesQueue) ReadFirst(queue string) (data interface{}){
    return m.queues[queue].Front().Value
    
}
func (m *DummyMessagesQueue) ReadFirstAll() (data map[string]interface{}){
    data = make(map[string]interface{});
    for key, value := range m.queues {
        data[key] = value.Front().Value;
    }
    return data
}
func NewDummyMessagesQueue(capacity int) (queue DummyMessagesQueue) {
    queue = DummyMessagesQueue{make(map[string]*list.List), capacity}
    return 
}
