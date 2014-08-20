package utils

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/ryadzenine/dolphin/models"
	"github.com/ryadzenine/dolphin/workers"
)

// Prepare The folds forthe V Fold
func PrepareFolds(nbFolds int, data models.SLDataset) map[int][2]int {
	out := make(map[int][2]int)
	fold_len := len(data) / nbFolds
	for i := 0; i < nbFolds; i++ {
		// We will build the content of each fold
		if i == nbFolds-1 {
			out[i] = [2]int{i * fold_len, len(data) - 1}
		} else {
			out[i] = [2]int{i * fold_len, (i + 1) * fold_len}
		}
	}
	return out
}

func FeedWorker(wg *sync.WaitGroup, data models.SLDataset,
	tdata []models.SLPoint, fold [2]int,
	worker *workers.Computable, worker_id int, total_workers int,
	flush chan<- bool) {
	defer wg.Done()
	for key, dt := range data {
		// cette données la est a moi et elle ne fait pas partie du training set
		if key%total_workers == worker_id && (key < fold[0] || key >= fold[1]) {
			worker.Input <- dt
			if key == 0 && flush != nil {
				flush <- true
			}
		}
	}
}

func MesureCompleteState(quit chan bool, workers []*workers.Computable, tdata []models.SLPoint) {
	ticker_c := 1
	ticker := time.Tick(30 * time.Second)
	for {
		select {
		case <-quit:
			return
		case <-ticker:
			// We will compute the L2 Errors
			go func() {
				states := make([]models.State, len(workers))
				for k, wk := range workers {
					states[k] = wk.Est.State()
				}
				av := models.States(states).Average()
				err := 0.0
				for key, v := range tdata {
					p := av[key]
					err = err + (p-v.Y)*(p-v.Y)
				}
				err = math.Sqrt(err) / float64(len(tdata))
				fmt.Println("Mes,", ticker_c*30, "Error,", err)
				ticker_c++
			}()
		default:
			continue
		}
	}
}

func MesureState(quit chan bool, worker *workers.Computable, tdata []models.SLPoint) {
	ticker_c := 1
	ticker := time.Tick(30 * time.Second)
	for {
		select {
		case <-quit:
			return
		case <-ticker:
			tt := time.Now()
			// We will compute the L2 Errors
			fmt.Println("L2:Mes:", worker.Id, ":", ticker_c*30, ":", worker.Est.Error(tdata), "time:", time.Since(tt))
			ticker_c++
		default:
			continue
		}
	}
}
