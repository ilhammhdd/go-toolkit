package goroutinekit

import (
	"os"
	"os/signal"
	"sync"
)

type Worker interface {
	Work() interface{}
}

type ResultHandler interface {
	Handle(interface{})
}

type Job interface {
	Worker
	ResultHandler
}

type WorkerPool struct {
	Job    chan Job
	Work   chan func()
	Worker chan Worker
	Done   chan bool
	PoolWG sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	wp := &WorkerPool{}

	wp.Work = make(chan func())
	wp.Job = make(chan Job)
	wp.Worker = make(chan Worker)
	wp.Done = make(chan bool)

	wp.PoolWG.Add(3)

	Do(func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

	JobLoop:
		for {
			select {
			case job := <-wp.Job:
				Do(func() {
					job.Handle(job.Work())
				})
			case <-wp.Done:
				break JobLoop
			case <-signals:
				break JobLoop
			}
		}

		close(wp.Job)
		wp.PoolWG.Done()
	})

	Do(func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

	WorkLoop:
		for {
			select {
			case work := <-wp.Work:
				Do(func() {
					work()
				})
			case <-wp.Done:
				break WorkLoop
			case <-signals:
				break WorkLoop
			}
		}

		close(wp.Work)
		wp.PoolWG.Done()
	})

	Do(func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

	WorkerLoop:
		for {
			select {
			case worker := <-wp.Worker:
				Do(func() {
					worker.Work()
				})
			case <-wp.Done:
				break WorkerLoop
			case <-signals:
				break WorkerLoop
			}
		}

		close(wp.Worker)
		wp.PoolWG.Done()
	})

	return wp
}
