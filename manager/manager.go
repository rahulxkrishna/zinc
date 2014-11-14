/*
Package manager implements the gron manager.
*/
package manager

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
	crontab string
	entries [MaxEntries]Entry
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

	fmt.Printf("Reading crontab file %s...\n", m.Crontab())

	buf, err := ioutil.ReadFile(m.Crontab())
	if err != nil {
		return x, err
	}
	fmt.Println(string(buf))

	return
}

// ServeHTTP is the HTTP server for gron
func (m *Manager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "gron!")
}
