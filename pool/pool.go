package pool

type Pool struct {
	sem   chan struct{}
	works chan func()
}

func New(size uint, queue uint) *Pool {
	p := &Pool{
		sem:   make(chan struct{}, size),
		works: make(chan func(), queue),
	}
	p.sem <- struct{}{}
	go p.worker(func() {})
	return p
}

func (p *Pool) Schedule(task func()) {
	p.schedule(task)
}

func (p *Pool) schedule(task func()) error {
	select {
	case p.sem <- struct{}{}:
		go p.worker(task)
		return nil
	case p.works <- task:
		return nil
	}
}

func (p *Pool) worker(task func()) {
	defer func() {
		<-p.sem
	}()
	task()
	for task := range p.works {
		task()
	}
}

func (p *Pool) Stop() {
	close(p.sem)
	close(p.works)
}
