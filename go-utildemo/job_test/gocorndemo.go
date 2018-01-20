package main

import (
	"fmt"

	"github.com/robfig/cron"
)

func main() {
	spec := "0, 5, 15, *, *, *"
	c := cron.New()
	c.AddFunc(spec, callFunc)
	c.Start()
	select {}
}

func callFunc() {
	fmt.Println("call come here")
}
