package eval

import (
	"fmt"
	"github.com/teacoder/gron/manager"
	_ "io/ioutil"
	_ "math/rand"
	"os/exec"
	"sort"
	_ "strings"
	"time"
)

type Exec struct {
	id       uint32
	atTime   int64
	interval int64 //seconds
	cmd      string
}

type byTime []Exec

var eQ [manager.MaxEntries]Exec

// XXX Temporary hack ! Got to re-organize the struct, store the
// queue count passed by the manager and return it.
func (e byTime) Len() int {
	count := 0
	for i := 0; i < len(e); i++ {
		if e[i].cmd != "" {
			count++
		}
	}
	return count
}
func (e byTime) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e byTime) Less(i, j int) bool { return e[i].atTime < e[j].atTime }

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
		eQ[i] = Exec{100, curTime.Unix() + secondsTd, int64(SecondsInDay), cmd}
		fmt.Printf("Running [%s] in %d seconds\n", cmd, secondsTd)
	}
	sort.Sort(byTime(eQ[0:]))
}

// Init initializes the internal data structures for the evaluator
func initialize(entries [manager.MaxEntries]manager.Entry, count int) {
	populateExecQueue(entries, count)
}

// runHead runs the commands which are up for execution starting at the head
// of the execution queue.
func runHead() {
	for i := 0; i < len(eQ); i++ {
		if eQ[i].id != 0 && time.Now().Unix() >= eQ[i].atTime {
			execute(eQ[i].cmd)
			//Now, schedule it for the next interval
			eQ[i].atTime += eQ[i].interval
			fmt.Printf("Executed [%s], next on %s \n",
				eQ[i].cmd,
				time.Unix(eQ[i].atTime, 0).String())
		} else {
			break
		}
	}

	sort.Sort(byTime(eQ[0:]))
}

// Run is the entry function into the evaluator. It initializes the execution queue
// and bocks till the time the earliest command is to be run.
func Run(entries [manager.MaxEntries]manager.Entry, count int) {
	initialize(entries, count)
	for {
		timer := time.NewTimer(time.Duration(eQ[0].atTime-time.Now().Unix()) * time.Second)
		<-timer.C
		fmt.Println("Here")
		runHead()
	}
}
