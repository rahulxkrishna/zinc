package main

import (
	"github.com/teacoder/gron/eval"
	"github.com/teacoder/gron/manager"
	//"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// create a new manager object
	man := manager.New()
	man.ReadCrontab()
	//http.ListenAndServe("localhost:40000", man)

	go eval.Run()

	wg.Wait()
}
