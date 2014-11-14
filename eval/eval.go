package eval

import (
	"fmt"
	"github.com/teacoder/gron/manager"
	_ "io/ioutil"
	_ "math/rand"
	"os/exec"
	_ "strings"
	"time"
)

type Exec struct {
	id       uint32
	atTime   int64
	interval int64 //seconds
	cmd      string
}

var execQueue [manager.MaxEntries]Exec

// Execute runs the command passed to it
func execute(cmd string) error {
	cmdExec := exec.Command(cmd)
	op, err := cmdExec.Output()

	if err != nil {
		fmt.Println("Error : " + err.Error())
		return err
	}

	fmt.Print(string(op))

	return nil
}

// PopulateExecQueue converts the in-memory config to an exec queue format for
// the evaluator to run over and execute.
func populateExecQueue() {

	// Let's deal with some dummy values now
	execQueue[0] = Exec{100, time.Now().Unix(), 10, "ls"}
	execQueue[1] = Exec{101, time.Now().Unix(), 5, "date"}
	execQueue[2] = Exec{102, time.Now().Unix(), 15, "time"}
}

// Init initializes the internal data structures for the evaluator
func initialize() {
	populateExecQueue()
}

// Poll will be called periodicaly from the main loop.
// It walks the execution queue and run any commands whose time has come
func poll() {
	for i := 0; i < len(execQueue); i++ {
		if execQueue[i].id != 0 && time.Now().Unix() >= execQueue[i].atTime {
			execute(execQueue[i].cmd)
			//Now, schedule it for the next interval
			execQueue[i].atTime += execQueue[i].interval
		}
	}
}

// Run is the entry function into the evaluator. It runs for ever, polling the
// execution queue every second
func Run() {
	fmt.Println("Running Evaluator")
	initialize()

	for {
		time.Sleep(1000)
		poll()
	}
}
