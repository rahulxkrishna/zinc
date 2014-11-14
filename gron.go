package main

import (
	"github.com/teacoder/gron/eval"
	"github.com/teacoder/gron/manager"
	//"net/http"
)

func main() {
	// create a new manager object
	man := manager.New()
	man.ReadCrontab()
	//http.ListenAndServe("localhost:40000", man)

	eval.Run()
}
