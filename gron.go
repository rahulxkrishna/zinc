package main

import (
	"github.com/teacoder/gron/eval"
	"github.com/teacoder/gron/manager"
)

func main() {
	// create a new manager object
	man := manager.New()
	man.ReadCrontab()

	eval.Run()
}
