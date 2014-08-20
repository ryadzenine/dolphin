package workers

import (
	"strconv"

	"github.com/ryadzenine/dolphin/models"
	"github.com/ryadzenine/dolphin/mpi"
)

// Represents a computable Work to be carried by a
// Worker
// TODO : Figure out how to make it an interface
type Computable struct {
	Id    int // TODO change Id to ID as lint says
	Name  string
	Input chan models.SLPoint
	// TODO Change to Regression Estimator
	Est models.Estimate
}

// Builds a new Computable
// queue: The Message Queue to be used by the worker
// id : the id of the computable
// points: The points to build the Estimator
// smooth: The smoothing parameter of the kernel estimator
func NewComputable(queue mpi.MessagesQueue, id int,
	est models.Estimate) *Computable {
	workerName := string(strconv.AppendInt([]byte("Worker "), int64(id), 10))
	ch := make(chan models.SLPoint)
	queue.Register(workerName)
	queue.Write(workerName, est.State())
	return &Computable{id, workerName, ch, est}
}
