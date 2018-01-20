package main

import (
	"fmt"
	"time"

	"github.com/panenming/go-im/libs/jobs"
	"github.com/panenming/go-im/libs/jobs/expression/every"
)

func main() {
	expr := every.NewExpression(10 * time.Second)
	jobs.Schedule(func() {
		fmt.Println("运行中。。。")
	}, expr)

	jobs.Start()
	time.Sleep(60 * time.Second)
	// 使用for{}堵塞会导致cpu飘高，导致cpu一个核100%
}
