package lib

type Task interface {
	Execute()
}

type Pool struct {
	// mu    sync.Mutex
	// size  int
	tasks chan Task
	kill  chan bool
	// wg    sync.WaitGroup
}

func NewPool(size, queue int) *Pool {
	pool := &Pool{
		tasks: make(chan Task, queue),
		kill:  make(chan bool),
	}
	// pool.Resize(size)

	for i := 0; i < size; i++ {
		go pool.worker()
	}

	return pool
}

func (p *Pool) worker() {
	// defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute()
		case <-p.kill:
			return
		}
	}
}

// func (p *Pool) Resize(n int) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	for p.size < n {
// 		p.size++
// 		p.wg.Add(1)
// 		go p.worker()
// 	}
// 	for p.size > n {
// 		p.size--
// 		p.kill <- struct{}{}
// 	}
// }

func (p *Pool) Close() {
	close(p.tasks)
}

// func (p *Pool) Wait() {
// 	p.wg.Wait()
// }

func (p *Pool) Exec(task Task) {
	p.tasks <- task
}
