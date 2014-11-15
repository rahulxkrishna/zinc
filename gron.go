package main

import (
	"github.com/teacoder/gron/eval"
	"github.com/teacoder/gron/manager"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Create a new manager object
	_, err := manager.New()
	if err != nil {
		log.Fatal(err)
	}

	// @teacoder: eval needs to use the manager object returned above so that
	// it can access the crontab entries. perhaps do that during
	// initialization (outside Run())? like: eval.initialize(), then go
	// eval.Run()?
	go eval.Run()

	wg.Wait()
}
