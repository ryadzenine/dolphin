package main

import (
    "github.com/ryadzenine/dolphin/mpi"
    "github.com/ryadzenine/dolphin/estimators"
    "flag"
    "time"
    "runtime"
    "fmt"
    "math/rand"
)
var workers = flag.Int("workers", 2, 
    "define how many workers will be launched")
func Worker(data_stream chan []float64,est *estimators.RevezEstimator ,queue mpi.MessagesQueue, tau int, name string){
    i := 0
    state := queue.ReadFirstAll()
    r:= <-data_stream
    queue.Write(name, r)
    for{
        data := <- data_stream 
        // ici on va faire des computations 
        if i % tau == 0 {
            tmp_state := queue.ReadFirstAll()
            count := 0
            acc := 0.0
            for i:=0; i < *workers ; i++ {
                key := string(i)
                if tmp_state[key] == state[key] {
                    tmp_state[key] = 0
                }else{
                    state[key] = tmp_state[key]
                    count = count + 1
                    acc = acc + state[key].(float64)
                }
            }
            acc = acc / float64(count) 
            est.ComputeDistributedStep(acc, data, 0 )
        }else{
            est.ComputeStep(data)
        }
        queue.Write(name, r)
        i=i+1
    }
}

func main(){
    flag.Parse()    
    runtime.GOMAXPROCS(runtime.NumCPU())
    queue := mpi.NewDummyMessagesQueue(5)
    var chans []chan []float64 = make([]chan []float64, *workers ) 
    // now we will launch the workers 
    for i:=0; i < *workers ; i ++ {
        chans[i] = make(chan []float64)
        est := estimators.NewRevezEstimator()
        go Worker(chans[i], est, &queue,10, "chan" + string(i))
    }
    rand.Seed(time.Now().UnixNano())
    for j:=0; j < 10000; j++ {
        for k:=0; k < *workers; k++ {
            chans[k] <- []float64{rand.Float64()}
        }
    }
}
