package workerpool

import (
	"sync"
)

type WorkerPool struct {
	jobs       chan Job
	workers    []Worker
	maxWorkers int
	stopOnce   sync.Once
	wg         sync.WaitGroup
}

func New(maxWorkers int) *WorkerPool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	pool := WorkerPool{
		jobs:       make(chan Job),
		workers:    make([]Worker, maxWorkers),
		maxWorkers: maxWorkers,
	}

	pool.start()

	return &pool
}

func (wp *WorkerPool) Submit(job Job) {
	if job == nil {
		return
	}
	wp.jobs <- job
}

func (wp *WorkerPool) Stop() {
	wp.stopOnce.Do(func() {
		close(wp.jobs)
		for _, worker := range wp.workers {
			worker.Stop()
		}
		wp.wg.Wait()
	})
}

func (wp *WorkerPool) start() {
	for i := 0; i < wp.maxWorkers; i++ {
		wp.wg.Add(1)

		worker := Worker{
			Jobs: wp.jobs,
			Quit: make(chan struct{}),
			wg:   &wp.wg,
		}

		worker.Run()

		wp.workers = append(wp.workers, worker)
	}
}

type Job func()

type Worker struct {
	Jobs chan Job
	Quit chan struct{}
	wg   *sync.WaitGroup
}

func NewWorker(jobs chan Job) Worker {
	return Worker{
		Jobs: jobs,
		Quit: make(chan struct{}),
	}
}

func (w Worker) Run() {
	go func() {
		for {
			select {
			case job, ok := <-w.Jobs:
				if ok {
					job()
				}
			case <-w.Quit:
				if w.wg != nil {
					w.wg.Done()
				}
				close(w.Quit)
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.Quit <- struct{}{}
	}()
}
