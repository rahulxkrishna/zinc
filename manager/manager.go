/*
Package manager implements the gron manager.
*/
package manager

import (
	//"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	//"os"
)

// Entry represents one entry in the crontab file
type Entry struct {
	minute uint8
	hour   uint8
	dom    uint8
	mon    uint8
	dow    uint8
	user   string
	cmd    string
}

// DefCrontab is the default crontab file
const DefCrontab = "/etc/crontab"

// MaxEntries is the maximum supported number of entries in the crontab
// file
const MaxEntries = 10

// Manager defines the manager struct
type Manager struct {
	crontab         string
	crontabContents string
	entries         [MaxEntries]Entry
}

// Crontab is the getter function for crontab
func (m *Manager) Crontab() string {
	return m.crontab
}

// SetCrontab sets the crontab path
func (m *Manager) SetCrontab(p string) {
	m.crontab = p
}

// New Creates a new manager
func New() *Manager {
	return &Manager{crontab: DefCrontab}
}

// ReadCrontab reads the crontab
func (m *Manager) ReadCrontab() (x int, err error) {
	x, err = 0, nil

	fmt.Printf("Reading crontab file %s...\n", m.crontab)
	buf, err := ioutil.ReadFile(m.crontab)
	if err != nil {
		return x, err
	}
	m.crontabContents = string(buf)

	/*
		f, err := os.Open(m.crontab)
		if err != nil {
			return x, err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		f.Close()
	*/

	return
}

// ServeHTTP is the HTTP server for gron
func (m *Manager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, m.crontabContents)
}
