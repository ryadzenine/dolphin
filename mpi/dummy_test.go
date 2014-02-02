package mpi

import "testing"
import "container/list"

func TestDummyMessagesQueues(t *testing.T) {
    mpi:= DummyMessagesQueue{ make(map[string]*list.List), 10}
    mpi.Register("q1", 5)
    if name,ok := mpi.queues["q1"]; ok {
        if name.Front().Value.(int) != 5 {
            t.Error("Expected an element in the DummyQueue")
        }

    }
    if mpi.ReadFirst("q1").(int) != 5 {
        t.Error("Expected Value 5")
    }
    mpi.Register("q2", 6)
    data := mpi.ReadFirstAll()
    if v, ok := data["q1"]; !ok || v.(int) != 5{
        t.Error("Expected queue Q1 to be filled with 5")
    }
    if v, ok := data["q2"]; !ok || v.(int) != 6{
        t.Error("Expected queue Q2 to be filled with 6")
    }
}

