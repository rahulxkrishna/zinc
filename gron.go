package main

import (
	"github.com/teacoder/gron/eval"
	"github.com/teacoder/gron/manager"
	"log"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// create a new manager object
	man := manager.New()
	_, err := man.ReadCrontab()
	if err != nil {
		log.Fatal(err)
	}
	go http.ListenAndServe("localhost:40000", man)

	go eval.Run()

	wg.Wait()
}
