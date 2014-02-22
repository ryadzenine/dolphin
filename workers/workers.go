package workers;

import (
    "github.com/ryadzenine/dolphin/mpi"
    "github.com/ryadzenine/dolphin/models/np"
)

func SimpleWorker(data_stream chan np.LearningPoint, est *np.RevezEstimator,
  queue mpi.MessagesQueue, tau int, name string){
    i := 1
    vc := make(map[string]int) // version control map 
    for{
        select{
        case data:= <- data_stream :
            if i == 1 {
                for _,v := range queue.Queues(){
                    vc[v] = 0
                }
            }
            // ici on va faire des computations 
            if i % tau == 0 {
                stat := queue.ReadStates(vc)
                // Block of code just to covert to the good types 
                states:= make([]np.State, 0, len(stat))
                for _,v := range stat {
                    states = append(states, v.(np.State))
                }
                // We know append the knew versions 
                for key,v := range stat {
                    vc[key] = v.Version()
                }
                acc := make([]float64, len(est.Points))
                if len(states) != 0 {
                    acc = np.States(states).ComputeAgregation()
                }
                est.ComputeDistributedStep(acc, data)
            }else{
                est.ComputeStep(data)
            }
            if i % tau > 4 {
               queue.Write(name, est.State())
           }
            i=i+1
        }
    }
}
