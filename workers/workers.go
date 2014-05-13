package workers

import (
	"github.com/ryadzenine/dolphin/models"
	"github.com/ryadzenine/dolphin/mpi"
)

func SimpleWorker(queue mpi.MessagesQueue, cmpt *Computable, tau int) {
	i := 1
	vc := make(map[string]int) // version control map
	for {
		select {
		case data, ok := <-cmpt.Input:
			if !ok {
				return
			}
			if i == 1 {
				for _, v := range queue.Queues() {
					vc[v] = 0
				}
			}
			// ici on va faire des computations
			if i%tau == 0 && tau != 1 {
				stat := queue.ReadStates(vc)
				// Block of code just to covert to the good types
				states := make([]models.State, 0, len(stat))
				for _, v := range stat {
					states = append(states, v.(models.State))
				}
				// We know append the knew versions
				for key, v := range stat {
					vc[key] = v.Version()
				}
				acc := make([]float64, len(cmpt.Est.Points))
				if len(states) != 0 {
					acc = models.States(states).ComputeAgregation()
				}
				cmpt.Est.ComputeDistributedStep(acc, data)
			} else {
				cmpt.Est.ComputeStep(data)
			}
			queue.Write(cmpt.Name, cmpt.Est.State())
			i = i + 1
		}
	}
}
