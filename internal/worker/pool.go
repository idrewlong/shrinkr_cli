package worker

import (
	"sync"

	"github.com/idrewlong/shrinkr_cli/internal/compressor"
)

// Pool manages concurrent image compression workers using goroutines.
// This is the key performance improvement over the Node.js version,
// which processes images sequentially.
type Pool struct {
	workerCount int
	jobs        chan compressor.Job
	results     chan compressor.Result
	wg          sync.WaitGroup
}

// NewPool creates a pool with the specified number of workers.
func NewPool(workerCount, totalJobs int) *Pool {
	// Buffer channels to prevent blocking
	bufSize := workerCount * 2
	if bufSize > totalJobs {
		bufSize = totalJobs
	}
	if bufSize < 1 {
		bufSize = 1
	}

	return &Pool{
		workerCount: workerCount,
		jobs:        make(chan compressor.Job, bufSize),
		results:     make(chan compressor.Result, bufSize),
	}
}

// Start launches all worker goroutines. Each reads from the jobs channel
// and sends results to the results channel.
func (p *Pool) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.jobs {
				result := compressor.Compress(job)
				p.results <- result
			}
		}()
	}

	// Close results channel when all workers are done
	go func() {
		p.wg.Wait()
		close(p.results)
	}()
}

// Submit adds a job to the queue.
func (p *Pool) Submit(job compressor.Job) {
	p.jobs <- job
}

// Done signals that no more jobs will be submitted.
func (p *Pool) Done() {
	close(p.jobs)
}

// Results returns the results channel for reading.
func (p *Pool) Results() <-chan compressor.Result {
	return p.results
}
