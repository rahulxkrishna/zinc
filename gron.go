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

	// @jayanthc : We should move the "Entry", and other common structures out of
	// the manager package and into a common 'gron' realm. It'd be nice to keep these
	// two loosely coupled.

	// Create a new manager object
	mgr, err := manager.New()
	if err != nil {
		log.Fatal(err)
	}

	go eval.Run(mgr.Entries(), int(mgr.EntryCount()))

	wg.Wait()
}
