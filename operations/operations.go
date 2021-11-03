package operations

import (
	"fmt"
	"log"
	"sync/atomic"
)

var (
	operations       = make(map[int64]Runnable)
	index      int64 = 0
	// Two channels - to send them work and to collect their results.
	jobs    = make(chan Runnable)
	results = make(chan Runnable)
)

type operationBase struct {
	status   Status
	loglines []string
	id       int64
}

type Status string

var (
	StatusError    = Status("error")
	StatusRunning  = Status("running")
	StatusFinished = Status("finished")
)

type Runnable interface {
	Run() error
	Log() []string
	Status() Status
	SetStatus(Status)
	SetID(int64)
	ID() int64
}

func Submit(operation Runnable) (int64, error) {
	id := nextIndex()
	operations[id] = operation
	jobs <- operation
	return id, nil
}

func GetStatus(id int64) (Status, error) {
	operation, ok := operations[id]
	if !ok {
		return "", fmt.Errorf("operation with id %d not found", id)
	}
	return operation.Status(), nil
}
func GetLogs(id int64) ([]string, error) {
	operation, ok := operations[id]
	if !ok {
		return []string{}, fmt.Errorf("operation with id %d not found", id)
	}
	return operation.Log(), nil
}

func WorkIt() {
	// This starts up 3 workers, initially blocked because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	log.Println("waiting for results")
	for {
		result := <-results
		log.Printf("Operation %d finished with result %s", result.ID(), result.Status())
	}
}

// this won't work yet, because nothing is consuming results
// maybe that channel doesn't need to exist since there's an ID
// on the Operation
func worker(id int, jobs <-chan Runnable, results chan<- Runnable) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		err := j.Run()
		if err != nil {
			fmt.Println("worker", id, "error on job", j.ID())
		} else {
			fmt.Println("worker", id, "finished", j)

		}

		results <- j
	}
}

func nextIndex() int64 {
	atomic.AddInt64(&index, 1)
	return index
}
