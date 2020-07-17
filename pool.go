package pool

// WorkFunc is the worker function user define
type WorkFunc func(arg ...interface{}) interface{}

type Pool struct {
	size    int
	blocker chan struct{}
}

// SetSize change the pool size
func (p *Pool) SetSize(size int) {
	p.size = size
	p.blocker = make(chan struct{}, size)
}

// Spawn do a job in pool
func (p *Pool) Spawn(f WorkFunc, args ...interface{}) *Job {
	p.blocker <- struct{}{}

	var job Job
	job.running()

	go func(f WorkFunc, args ...interface{}) {
		job.start(f, args...)
		<-p.blocker
	}(f, args...)

	return &job
}

// JoinAll wait all job finished
func (p *Pool) JoinAll(jobs []*Job) {
	for _, job := range jobs {
		job.Join()
	}
}

// Map do multiple jobs and return results in order
func (p *Pool) Map(f WorkFunc, args []interface{}) []interface{} {
	if len(args) == 0 {
		return []interface{}{}
	}

	jobs := make([]*Job, len(args))
	for i, arg := range args {
		job := p.Spawn(f, arg)
		jobs[i] = job
	}

	result := make([]interface{}, len(args))
	for i, job := range jobs {
		result[i] = job.Get()
	}

	return result
}

// New make a new Pool
func New(size int) *Pool {
	return &Pool{
		size:    size,
		blocker: make(chan struct{}, size),
	}
}
