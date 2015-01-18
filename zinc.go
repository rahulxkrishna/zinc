package main

import (
	"github.com/1d4Nf6/zinc/eval"
	"github.com/1d4Nf6/zinc/manager"
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
