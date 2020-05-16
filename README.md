## workerpool

[![Go Report Card](https://goreportcard.com/badge/github.com/ostapenkoden/workerpool)](https://goreportcard.com/report/github.com/ostapenkoden/workerpool)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/ostapenkoden/workerpool/blob/master/LICENSE)

Simple worker pool implementation for concurrency control in Golang.
Limits the concurrency of job execution, not the number of jobs queued. 
Never blocks submitting jobs, no matter how many jobs are queued.

## Installation
```
$ go get github.com/ostapenkoden/workerpool
```

## Example
```go
package main

import (
	"fmt"
	"github.com/ostapenkoden/workerpool"
)

func main() {
	wp := workerpool.New(3)
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, job := range jobs {
		job := job
		wp.Submit(func(){
			fmt.Printf("Handling job %d\n", job)
		})
	}
	wp.Stop()
}
```