package lib

// Task -
type Task interface {
	Execute(id int)
}

// Pool -
type Pool struct {
	// mu    sync.Mutex
	// size  int
	tasks chan Task
	kill  chan bool
	// wg    sync.WaitGroup
}

// NewPool -
func NewPool(size, queue int) *Pool {
	pool := &Pool{
		tasks: make(chan Task, queue),
		kill:  make(chan bool),
	}
	// pool.Resize(size)

	for i := 0; i < size; i++ {
		go pool.worker(i)
	}

	return pool
}

func (p *Pool) worker(id int) {
	// defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute(id)
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

// Close -
func (p *Pool) Close() {
	close(p.tasks)
}

// func (p *Pool) Wait() {
// 	p.wg.Wait()
// }

// Exec -
func (p *Pool) Exec(task Task) {
	p.tasks <- task
}
