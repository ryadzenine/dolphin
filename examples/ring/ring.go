package main

import (
	"encoding/gob"
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

// The ring configuration file a json containing the list of hosts
// that will take part in the computations
var network = flag.String("network", "", "a json file with the adresses of the nodes of the ring")

// The position of the actual process in the ring
var me = flag.Int("me", 0, "my position in the network")

func VFoldCv(ring *Ring, queue *mpi.CircularMPI, data []models.SLPoint, len_ring int) {
	// We compute the fold size
	folds := exutils.PrepareFolds(5, data)
	for k, fold := range folds {
		fmt.Println("INFO:", *me, " Entering fold ", k)
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
			cp := workers.NewComputable(queue, *me*10+i, est)
			cmpt = append(cmpt, cp)
		}
		for i := 0; i < *workrs; i++ {
			go workers.SimpleWorker(queue, cmpt[i], *tau)
		}

		var wg sync.WaitGroup
		wg.Add(*workrs)
		totalWorkers := len_ring * (*workrs)
		j := 0
		fl := make(chan bool)
		for i := *workrs * (*me - 1); i < *workrs*(*me); i++ {
			if j == 0 && *me == 1 {
				go exutils.FeedWorker(&wg, data, tdata, fold, cmpt[j], i, totalWorkers, fl)
			} else {
				go exutils.FeedWorker(&wg, data, tdata, fold, cmpt[j], i, totalWorkers, nil)
			}
			j++
		}
		if *me == 1 {
			go func() {
				<-fl
				time.Sleep(250 * time.Millisecond)
				err := queue.Flush()
				if err != nil {
					fmt.Println("ERROR: Can't Flush the Queue, ", err)
				}
			}()
		}
		quit := make(chan bool)
		go exutils.MesureState(quit, cmpt[0], tdata)
		wg.Wait()
		fmt.Println("L2:FINAL:", cmpt[0].Id, ":", cmpt[0].Est.Error(tdata))
		quit <- true
		return
	}
}
func main() {
	runtime.GOMAXPROCS(8)
	flag.Parse()
	if *workrs > 6 || *workrs < 1 {
		panic("Wrong number of workers")
	}
	gob.Register(np.EstimatorState{})
	ring, err := NewRing(*me, *network)

	if err != nil {
		panic(fmt.Sprint("Ring Error:", err))
	}
	queue := mpi.NewCircularMPI("tcp", ring.Me, ring.Hosts, buildLogger("ring", *me))
	// First if we are the number one of the cluster
	// we nee to wait for the others.
	ring.Sync()

	// Now we got all of them so we can begin
	// We First Launch the Message Handler for
	go queue.ListenAndServe()
	learning, err := ioutil.ReadFile(*ldata)
	data := utils.ParseData(learning)
	top := time.Now()
	VFoldCv(ring, queue, data, len(ring.Hosts))
	finish := time.Since(top)
	fmt.Println("TIME:FINAL:", *me, finish)
}
