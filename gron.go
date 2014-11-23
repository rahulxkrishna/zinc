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

	// @jayanthc : We should move the "Entry", and other common structures out of
	// the manager package and into a common 'gron' realm. It'd be nice to keep these
	// two loosely coupled.

	// @teacoder: My idea was that eval would keep a reference to mgr that
	// would enable it to get all entries using mgr.Entries() instead of having
	// a special Get() method for each entry. From a design pov, I think it's
	// the manager's job to keep the entries, and the evaluator's job to ask
	// for it, rather than keep entries outside. That will also keep gron.go
	// simple.
	// PS: We should use issues for this sort of thing.

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
