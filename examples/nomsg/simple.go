package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"
	"time"

	exutils "github.com/ryadzenine/dolphin/examples/utils"
	"github.com/ryadzenine/dolphin/models"
	"github.com/ryadzenine/dolphin/models/np"
	"github.com/ryadzenine/dolphin/mpi"
	"github.com/ryadzenine/dolphin/utils"
	"github.com/ryadzenine/dolphin/workers"
)

// The number of workers that will handle de the computations
// two is juste fine
var workrs = flag.Int("workers", 2,
	"define how many workers will be launched")

// The value of tau determines the number of computations to be done
// before doing an aggregation
var tau = flag.Int("tau", 2,
	"tau defines the numbers of steps that have to bo computed by each worker before an agregation")

// Learning data formated as follow y;x1;x2;x3.... in a csv file
var ldata = flag.String("learning-data", "", "the learning dataset")

func VFoldCv(queue mpi.Dummy, data []models.SLPoint) {
	// We compute the fold size
	folds := exutils.PrepareFolds(5, data)
	for k, fold := range folds {
		fmt.Println("Entering fold ", k)
		// Now we build the training dataset for this fold
		tdata := data[fold[0]:fold[1]]
		tpoints := make([]models.Vector, 0, len(tdata))
		for _, pt := range tdata {
			tpoints = append(tpoints, models.Vector(pt.X))
		}

		// Now We will build the Workers
		cmpt := make([]*workers.Computable, 0, *workrs)
		for i := 0; i < *workrs; i++ {

			est, _ := np.NewRevezEstimator(tpoints)
			cp := workers.NewComputable(queue, i, est)
			cmpt = append(cmpt, cp)
		}
		for i := 0; i < *workrs; i++ {
			go workers.SimpleWorker(queue, cmpt[i], *tau)
		}

		var wg sync.WaitGroup
		wg.Add(*workrs)
		for i := 0; i < *workrs; i++ {
			go func(j int) {
				for k, val := range data {
					if k < fold[0] || k > fold[1] {
						cmpt[j].Input <- val
					}
				}
				wg.Done()
			}(i)
		}
		quit := make(chan bool)
		go exutils.MesureCompleteState(quit, cmpt, tdata)
		wg.Wait()
		fmt.Println("Final;", cmpt[0].Est.Error(tdata))
		quit <- true
	}
}
func main() {
	runtime.GOMAXPROCS(8)
	flag.Parse()
	queue := mpi.NewDummy()

	learning, err := ioutil.ReadFile(*ldata)
	if err != nil {
		panic(err)
		return
	}
	data := utils.ParseData(learning)
	top := time.Now()
	VFoldCv(queue, data)
	finish := time.Since(top)
	fmt.Println("Time;", finish)
}
