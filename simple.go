package main

import (
    "flag"
    "runtime"
    "strings"
    "bufio"
    "fmt"
    "os"
    "strconv"
    "io/ioutil"
    "github.com/ryadzenine/dolphin/mpi"
    "github.com/ryadzenine/dolphin/estimators"
    "github.com/ryadzenine/dolphin/workers"
)
// The number of workers that will handle de the computations
// usually as much as the number of thread (or cores) of you CPU
var workrs = flag.Int("workers", 2, 
    "define how many workers will be launched")
var tau = flag.Int("tau", 2, 
    "tau defines the numbers of steps that have to bo computed by each worker before an agregation")
// The File on wich we are going to learn the data + the number of data points 
// to use for the learning the rest will be user for the cross validation
var ldata = flag.String("learning-data", "", "the learning dataset")
var tdata = flag.String("test-data", "", "the testing dataset to compute to L2 empirical error")
// The interval on wich to do the estimation of Y 
var dimension = flag.Int("dim", 1 , "the dimension of the learning data points ")
var lowerBound = flag.Float64("inf", 0 , "The lower bound of the interval on wich we estimate Y")
var upperBound = flag.Float64("sup", 1 , "The upper bound of the interval on wich we estimate Y")
var domainStep = flag.Float64("step", 0.1, "the step of the cubic subdivision of (inf, sup)^d")
var smoothing = flag.Float64("smooth", 0.1, "the smoothing")


func StreamDataFromFile(step int, sc *bufio.Scanner, chans []chan estimators.LearningPoint){
    i := 0
    mod := len(chans)
    tot := 1 
    for sc.Scan() && tot <= step {
       i = (i+1) % mod  
       data :=  estimators.ParseLearningPoint(sc.Text())
       tot++ 
       chans[i] <- data
    }
}
func main(){
    flag.Parse()    
    runtime.GOMAXPROCS(runtime.NumCPU())
    // We create here a queue for the messaging 
    queue := mpi.NewDummyMessagesQueue(5)
    ests := make([] *estimators.RevezEstimator, *workrs)
    chans := make([]chan estimators.LearningPoint, *workrs ) 

    points, err := estimators.MeshEvalPoints(*lowerBound, *upperBound, *domainStep, *dimension)
    if err != nil {
        panic(err)
    }
    // now we will launch the workers 
    for i:=0; i < *workrs ; i ++ {
        chans[i] = make(chan estimators.LearningPoint)
        est,_ := estimators.NewRevezEstimator(points, *smoothing)
        ests[i] = est
        worker_name :=  string(strconv.AppendInt([]byte("name ") ,int64(i), 10 ))
        queue.Register(worker_name)
        queue.Write(worker_name, est.State())
    }
    for i:=0; i< *workrs; i++ {
        worker_name :=  string(strconv.AppendInt([]byte("name ") ,int64(i), 10 ))
        go workers.SimpleWorker(chans[i], ests[i], &queue, *tau, worker_name)
    }
    learning, err := os.Open(*ldata)
    if err != nil {
        panic(err)
    }
    sc := bufio.NewScanner(learning)
    test, err := ioutil.ReadFile(*tdata)
    if err != nil {
        panic(err) 
    }
    testData := make([]estimators.LearningPoint, 0,*dimension*len(test))
    for _,s := range strings.Split(string(test), "\n") {
        testData = append(testData, estimators.ParseLearningPoint(s))
    }
    for i:= 1 ; i <= 100; i++{
        StreamDataFromFile(200, sc, chans)
        err := ests[0].L2Error(testData[0:1000])
        fmt.Print(i*200)
        fmt.Print(";")
        fmt.Println(err)
    }
}
