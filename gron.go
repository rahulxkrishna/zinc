package main

import (
	"github.com/teacoder/gron/eval"
	"github.com/teacoder/gron/manager"
)

func main() {
	manager.ReadCrontab("gron.go")
	eval.Init()
}
