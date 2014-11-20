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
func populateExecQueue(entries [manager.MaxEntries]manager.Entry, count int) {

	const SecondsInMin int = 60
	const SecondsInHour int = 3600
	const SecondsInDay int = 86400
	var secondsTd int64

	///XXX Let's not worry about the DOM, DOW etc for now.
	curTime := time.Now()

	// There has to be a better way to do this.
	// TODO Add UTs
	secondsToday := SecondsInHour*curTime.Hour() + SecondsInMin*curTime.Minute()

	for i := 0; i < count; i++ {
		e := entries[i].Get()
		min := int(e["min"].(uint8))
		hr := int(e["hr"].(uint8))
		cmd := e["cmd"].(string)
		secondsSched := SecondsInHour*hr + SecondsInMin*min

		// Scheduled in the future on the same day
		if secondsSched > secondsToday {
			secondsTd = int64(secondsSched - secondsToday)
		} else { // Scheduled for the next day
			secondsTd = int64(((SecondsInHour * 24) - secondsToday) + secondsSched)
		}
		execQueue[i] = Exec{100, curTime.Unix() + secondsTd, int64(SecondsInDay), cmd}
		fmt.Printf("Running [%s] in %d seconds\n", cmd, secondsTd)
	}
}

// Init initializes the internal data structures for the evaluator
func initialize(entries [manager.MaxEntries]manager.Entry, count int) {
	populateExecQueue(entries, count)
}

// Poll will be called periodicaly from the main loop.
// It walks the execution queue and run any commands whose time has come
func poll() {
	for i := 0; i < len(execQueue); i++ {
		if execQueue[i].id != 0 && time.Now().Unix() >= execQueue[i].atTime {
			execute(execQueue[i].cmd)
			//Now, schedule it for the next interval
			execQueue[i].atTime += execQueue[i].interval
			fmt.Printf("Executed [%s], next on %s \n",
				execQueue[i].cmd,
				time.Unix(execQueue[i].atTime, 0).String())
		}
	}
}

// Run is the entry function into the evaluator. It runs for ever, polling the
// execution queue every second
func Run(entries [manager.MaxEntries]manager.Entry, count int) {
	initialize(entries, count)

	for {
		time.Sleep(1000)
		poll()
	}
}
