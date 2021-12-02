package pool

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// WorkFunc is the worker function user define
type WorkFunc func(arg ...interface{}) interface{}

type Pool struct {
	size  int64
	limit *semaphore.Weighted
}

// SetSize change the pool size
func (p *Pool) SetSize(size int64) {
	p.size = size
	p.limit = semaphore.NewWeighted(size)
}

// GetSize get the pool size
func (p *Pool) GetSize() int64 {
	return p.size
}

// Spawn do a job in pool
func (p *Pool) Spawn(f WorkFunc, args ...interface{}) *Job {
	limit := p.limit
	limit.Acquire(context.Background(), 1)

	var job Job
	job.running()

	go func(f WorkFunc, args ...interface{}) {
		job.start(f, args...)
		limit.Release(1)
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
func New(size int64) *Pool {
	return &Pool{
		size:  size,
		limit: semaphore.NewWeighted(size),
	}
}
