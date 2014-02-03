package mpi
//This interface represent a Message passing interface, it serve essentialy as 

// an abstraction to messageQueues or caching technologies
type MessagesQueue interface {
    // Register a message in the queue 
    Write(string, interface {})
    // Read the first message of 
    ReadFirst(string) interface{}
    // Returns the last message registred in every queue
    ReadFirstAll() map[string]interface{}
}

