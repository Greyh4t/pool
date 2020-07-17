package pool

import "sync"

type Job struct {
	wg    sync.WaitGroup
	value interface{}
}

func (r *Job) start(f WorkFunc, arg ...interface{}) {
	r.value = f(arg...)
	r.wg.Done()
}

func (r *Job) running() {
	r.wg.Add(1)
}

// Join wait job finished
func (r *Job) Join() {
	r.wg.Wait()
}

// Get wait and return the job result
func (r *Job) Get() interface{} {
	r.Join()
	return r.value
}
