package workers

import (
	"github.com/ryadzenine/dolphin/models"
	"github.com/ryadzenine/dolphin/models/np"
	"github.com/ryadzenine/dolphin/mpi"
	"strconv"
)

// Represents a computable Work to be carried by a
// Worker
// TODO : Change *np.RevezEstimator to an intreface
// TODO : Figure out how to make it an interface
type Computable struct {
	Id    int
	Name  string
	Input chan models.SLPoint
	Est   *np.RevezEstimator
}

// Builds a new Computable
// queue: The Message Queue to be used by the worker
// id : the id of the computable
// points: The points to build the Estimator
// smooth: The smoothing parameter of the kernel estimator
func NewComputable(queue mpi.MessagesQueue, id int,
	points []models.Point) *Computable {

	worker_name := string(strconv.AppendInt([]byte("Worker "), int64(id), 10))
	ch := make(chan models.SLPoint)
	est, _ := np.NewRevezEstimator(points)
	queue.Register(workerName)
	queue.Write(workerName, est.State())
	return &Computable{id, workerName, ch, est}
}
