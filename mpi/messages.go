package mpi
//This interface represent a Message passing interface, it serve essentialy as 
type Versionable interface {
    Version() int 
}
// an abstraction to messageQueues or caching technologies
type MessagesQueue interface {
    //Returns the queues name that are already registered in the MPI 
    Queues() []string
    // Create a new queue
    Register(string) bool
    // Register a message in the queue 
    Write(string, Versionable)
    // Read the first message of 
    ReadFirst(string) Versionable
    // Returns the last message registred in every queue
    ReadFirstAll() map[string]Versionable
    // Returns the last messages registred in every queue discarding 
    // the old ones 
    ReadStates(map[string]int) map[string]Versionable
}

