package mapreduce

import (
	"fmt"
	"sync"
)

//
// schedule() starts and waits for all tasks in the given phase (mapPhase
// or reducePhase). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	//fmt.Printf("DEBUG: schedule starts runing\n")
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers. Once all tasks
	// have completed successfully, schedule() should return.
	//
	// Your code here (Part III, Part IV).
	//

	var wg sync.WaitGroup
	for i := 0; i < ntasks; i++ {
		wg.Add(1)
		go func(taskInd int, n_other int, phase jobPhase){
				for{
					workerAddr := <-registerChan
					//fmt.Printf("DEBUG: task %d is assigned to worker with address %s\n", taskInd, workerAddr)
					var args DoTaskArgs
					args.JobName = jobName
					args.File = mapFiles[taskInd]
					args.Phase = phase
					args.TaskNumber = taskInd
					args.NumOtherPhase = n_other
					ok := call(workerAddr, "Worker.DoTask", &args, nil)
					if ok {
						wg.Done()
						registerChan <- workerAddr
						break
					}
				}
			}(i, n_other, phase)
	}
	wg.Wait()
	fmt.Printf("Schedule: %v done\n", phase)
}
