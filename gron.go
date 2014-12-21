package main

import (
	"github.com/teacoder/gron/eval"
	"github.com/teacoder/gron/manager"
	"log"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	var crontab string
	if len(os.Args) == 2 {
		crontab = os.Args[1]
	}

	// Create a new manager object
	mgr, err := manager.New(crontab)
	if err != nil {
		log.Fatal(err)
	}

	go eval.Run(mgr.Entries(), int(mgr.EntryCount()))

	wg.Wait()
}
