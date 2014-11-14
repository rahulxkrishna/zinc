package main

import (
	"github.com/teacoder/gron/eval"
	"github.com/teacoder/gron/manager"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// create a new manager object
	man := manager.New()
	man.ReadCrontab()

	go eval.Run()

	wg.Wait()
}
