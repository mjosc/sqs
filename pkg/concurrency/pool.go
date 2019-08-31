package concurrency

import (
	"sync"

	"github.com/mjosc/sqs/pkg/common"
)

func NewPool(nThreads int) *Pool {
	pool := Pool{
		nThreads: nThreads,
		tasks:    make(chan common.Task),
	}
	pool.build()
	return &pool
}

type Pool struct {
	nThreads int
	wg       sync.WaitGroup
	tasks    chan common.Task
}

func (p *Pool) Execute(task common.Task) {
	p.wg.Add(1)
	p.tasks <- task
}

func (p *Pool) build() {
	for i := 0; i < p.nThreads; i++ {
		go func() {
			for task := range p.tasks {
				defer p.wg.Done()
				task.Run()
			}
		}()
	}
}

func (p *Pool) Close() {
	close(p.tasks)
	p.wg.Wait()
}
