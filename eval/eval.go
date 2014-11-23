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

type ExecEntry struct {
	id       uint32
	atTime   int64
	interval int64 //seconds
	cmd      string
}

type byTime []ExecEntry

type EvalQueue struct {
	entries [manager.MaxEntries]ExecEntry
	count   int
}

type Evaluator struct {
	eQ EvalQueue
}

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

func (e byTime) Len() int           { return len(e) }
func (e byTime) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e byTime) Less(i, j int) bool { return e[i].atTime < e[j].atTime }

func (eQ *EvalQueue) sort() {
	sort.Sort(byTime(eQ.entries[0:eQ.count]))
}

// populateExecQueue converts the in-memory crontab config to an exec queue format for
// the evaluator to run over and execute.
func (e *Evaluator) populateExecQueue(
	entries [manager.MaxEntries]manager.Entry,
	count int) {

	const SecondsInMin int = 60
	const SecondsInHour int = 3600
	const SecondsInDay int = 86400
	var secondsTd int64

	///XXX Let's not worry about the DOM, DOW etc for now.
	curTime := time.Now()

	// There has to be a better way to do this.
	// TODO Add UTs
	secondsToday := SecondsInHour*curTime.Hour() + SecondsInMin*curTime.Minute() + curTime.Second()

	// @teacoder: Processing a single entry in a function will let us use the
	// unit-testing features of Go more effectively
	for i := 0; i < count; i++ {
		entry := entries[i].Get()
		min := int(entry["min"].(uint8))
		hr := int(entry["hr"].(uint8))
		cmd := entry["cmd"].(string)
		secondsSched := SecondsInHour*hr + SecondsInMin*min

		// Scheduled in the future on the same day
		if secondsSched > secondsToday {
			secondsTd = int64(secondsSched - secondsToday)
		} else { // Scheduled for the next day
			secondsTd = int64(((SecondsInHour * 24) - secondsToday) + secondsSched)
		}
		e.eQ.entries[i] = ExecEntry{100, curTime.Unix() + secondsTd, int64(SecondsInDay), cmd}
		fmt.Printf("Running [%s] in %d seconds\n", cmd, secondsTd)
	}
	e.eQ.count = count
	e.eQ.sort()
}

// initialize the internal data structures for the evaluator
func (e *Evaluator) initialize(entries [manager.MaxEntries]manager.Entry, count int) {
	e.populateExecQueue(entries, count)
}

// runHead runs the commands which are up for execution starting at the head
// of the execution queue.
func (e *Evaluator) runHead() {
	for i := 0; i < e.eQ.count; i++ {
		//If it's time, run it and schedule it for the next interval
		if time.Now().Unix() >= e.eQ.entries[i].atTime {
			execute(e.eQ.entries[i].cmd)
			e.eQ.entries[i].atTime += e.eQ.entries[i].interval
			fmt.Printf("Executed [%s], next on %s \n",
				e.eQ.entries[i].cmd,
				time.Unix(e.eQ.entries[i].atTime, 0).String())
		} else {
			break
		}
	}

	// Always keep the next event at the head of the queue
	e.eQ.sort()
}

// Run is the entry function into the evaluator. It initializes the execution queue
// and bocks till the time the earliest command is to be run.
func Run(entries [manager.MaxEntries]manager.Entry, count int) {
	var eval Evaluator
	eval.initialize(entries, count)
	for {
		timer := time.NewTimer(time.Duration(
			eval.eQ.entries[0].atTime-time.Now().Unix()) * time.Second)
		<-timer.C
		eval.runHead()
	}
}
