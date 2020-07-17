# pool
一个类似python gevent.Pool的协程池

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/greyh4t/pool"
)

func work(args ...interface{}) interface{} {
	fmt.Println(args...)
	time.Sleep(time.Second)
	return fmt.Sprintf("result:%v", args[0])
}

func example1() {
	p := pool.New(2)
	taskCount := 10

	var jobs []*pool.Job
	for i := 0; i < taskCount; i++ {
		job := p.Spawn(work, i)
		jobs = append(jobs, job)
	}

	// job.Get会自动等待任务结束
	for _, job := range jobs {
		result := job.Get()
		fmt.Println(result)
	}
}

func example2() {
	p := pool.New(2)
	taskCount := 10

	var jobs []*pool.Job
	for i := 0; i < taskCount; i++ {
		job := p.Spawn(work, i)
		jobs = append(jobs, job)
	}

	// 如果不需要获取结果，可以直接调用JoinAll等待任务结束
	p.JoinAll(jobs)
}

func example3() {
	p := pool.New(2)
	args := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// 调用Map直接获取所有参数的执行结果
	result := p.Map(work, args)

	fmt.Println(result)
}

func main() {
	fmt.Println("<--example1-->")
	example1()
	fmt.Println("<--example2-->")
	example2()
	fmt.Println("<--example3-->")
	example3()
}
```

## Installation

go get github.com/greyh4t/pool


