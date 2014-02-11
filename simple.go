package main

import (
    "github.com/ryadzenine/dolphin/mpi"
    "github.com/ryadzenine/dolphin/estimators"
    "flag"
    "time"
    "runtime"
    "fmt"
    "strconv"
    "math/rand"
)
var workers = flag.Int("workers", 2, 
    "define how many workers will be launched")
func Worker(data_stream chan estimators.LearningPoint,est *estimators.RevezEstimator ,queue mpi.MessagesQueue, tau int, name string){
    i := 1
    state := queue.ReadFirstAll()
    for{
        select{
        case data:= <- data_stream :
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
                est.ComputeDistributedStep(acc, data)
            }else{
                est.ComputeStep(data)
            }
            queue.Write(name, est.State)
            i=i+1
        }
    }
}

func main(){
    flag.Parse()    
    runtime.GOMAXPROCS(runtime.NumCPU())
    queue := mpi.NewDummyMessagesQueue(5)
    var chans []chan estimators.LearningPoint = make([]chan estimators.LearningPoint, *workers ) 
    // now we will launch the workers 
    for i:=0; i < *workers ; i ++ {
        chans[i] = make(chan estimators.LearningPoint)
        est := estimators.NewRevezEstimator()
        go Worker(chans[i], est, &queue,10, string(strconv.AppendInt([]byte("name ") ,int64(i), 10 )))
    }
    rand.Seed(time.Now().UnixNano())
    fmt.Println("Je m'apprete a envoyer du flux")
    for j:=0; j < 20; j++ {
        for k:=0; k < *workers; k++ {
            p := estimators.LearningPoint{ 
                    X: []float64{rand.Float64(), rand.Float64()}, 
                    Y: rand.Float64() }
            chans[k] <- p
        }
    }
}
