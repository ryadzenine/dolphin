package mpi

import "net/http"

// Versionalble represent a versionable object it is used by the Message Queue to sort the messages
// in a LIFO style.
type Versioner interface {
	Version() int
}

// MessageQueue an abstraction to messageQueues or caching technologies
type MessagesQueue interface {
	//Returns the queues name that are already registered in the MPI
	Queues() []string
	// Create a new queue
	Register(string) bool
	// Register a message in the queue
	Write(string, Versioner)
	// Read the first message of
	// TODO : Rajouter une gestion d'erreur
	ReadFirst(string) Versioner
	// Returns the last message registred in every queue
	ReadFirstAll() map[string]Versioner
	// Returns the last messages registred in every queue discarding
	// the old ones
	ReadStates(map[string]int) map[string]Versioner
}

// NMessQueue TODO see what to do with it
type NMessQueue interface {
	MessagesQueue
	LocalQueues() []string
	MessagesHandler(w http.ResponseWriter, r *http.Request)
}
