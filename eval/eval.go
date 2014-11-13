package eval

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

const CONFIG_FILE = "crontab"
const MAX_CRONTAB_SZ = 10

type Config struct {
	id         uint32
	minute     string
	hour       string
	dayOfMonth string
	month      string
	dayOfWeek  string
	action     string
}

var runningConfig [MAX_CRONTAB_SZ]Config

// RedConfig reads the configuration in when starting up
// Just a raw representation of the file data; will be converted to an
// appropriate internal DS.
// The crontab file is expected to be in the standard crontab format

//TODO Comments made by the user in the crontab file should not be lost,
// should we add/remove rules via our interface and then submit to the file.

func ReadConfig() {
	rand.Seed(time.Now().Unix())

	// Add a uid to identify a entry uniquely (Helpful for the manager to add/del?)
	uid := rand.Uint32()

	buf, err := ioutil.ReadFile(CONFIG_FILE)

	if err != nil {
		fmt.Println("Failed to open crontab, starting up with an empty config")
	}

	lines := strings.Split(string(buf), "\n")

	for i := 0; i < len(lines); i++ {
		attrs := strings.Fields(lines[i])
		if len(attrs) > 0 {
			runningConfig[i].id = uid
			runningConfig[i].minute = attrs[0]
			runningConfig[i].hour = attrs[1]
			runningConfig[i].dayOfMonth = attrs[2]
			runningConfig[i].month = attrs[3]
			runningConfig[i].dayOfWeek = attrs[4]
			runningConfig[i].action = attrs[5]
			uid += 1
		}
	}
}

// PrintConfig is a utility function to dump the in-memory config

func PrintConfig() {
	for i := 0; i < MAX_CRONTAB_SZ; i++ {
		if runningConfig[i].id > 0 {
			fmt.Printf("%d %s %s %s %s %s %s\n",
				runningConfig[i].id,
				runningConfig[i].minute,
				runningConfig[i].hour,
				runningConfig[i].dayOfMonth,
				runningConfig[i].month,
				runningConfig[i].dayOfWeek,
				runningConfig[i].action)
		}
	}
}

func Run() {
	fmt.Println("Running the Evaluator")
	ReadConfig()
	PrintConfig()
}
