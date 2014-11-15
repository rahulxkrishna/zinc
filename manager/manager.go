/*
Package manager implements the gron manager.
*/
package manager

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	//"reflect"
	"strconv"
	"strings"
)

// init initializes the manager package
func init() {
	// Check user, permissions?
}

// Entry represents one entry in the crontab file
type Entry struct {
	min  uint8
	hr   uint8
	dom  uint8
	mon  uint8
	dow  uint8
	user string
	cmd  string
}

// Default values
// defCrontab is the default crontab file
const defCrontab = "/etc/crontab"
const defHostname = "localhost"
const defPort = 40000

// MaxEntries is the maximum supported number of entries in the crontab file
const MaxEntries = 10

// Manager defines the manager struct
type Manager struct {
	// Crontab-related variables
	crontab         string
	crontabContents string
	entries         [MaxEntries]Entry
	entryCount      uint8 // TODO: think of a better data type
	// Webserver-related variables
	hostname string
	port     uint16
}

// New Creates a new manager and performs initialization
func New() (m *Manager, err error) {
	// Create and initialize manager
	m = &Manager{crontab: defCrontab, hostname: defHostname, port: defPort}
	err = m.readCrontab()

	// Run webserver in a goroutine
	go http.ListenAndServe(m.hostname+":"+strconv.Itoa(int(m.port)), m)

	return m, err
}

// Entries gets the crontab entries stored in the Manager
func (m *Manager) Entries() [MaxEntries]Entry {
	return m.entries
}

// EntryCount gets the number of crontab entries stored in the Manager
func (m *Manager) EntryCount() uint8 {
	return m.entryCount
}

// readCrontab reads the crontab and builds the entries data structure
func (m *Manager) readCrontab() (err error) {
	err = nil

	// Read the crontab
	fmt.Printf("Reading crontab file %s...\n", m.crontab)
	buf, err := ioutil.ReadFile(m.crontab)
	if err != nil {
		return
	}
	m.crontabContents = string(buf)

	// Scan the crontab lines
	// TODO: is there a simpler way to do this?
	scanner := bufio.NewScanner(strings.NewReader(m.crontabContents))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			// Empty line; Ignore
			continue
		}
		line = strings.Replace(line, "*", "0", -1)
		tokens := strings.Fields(line)
		if tokens[0] == "#" {
			// Comment line; Ignore
			continue
		} else if len(tokens) == 1 {
			// Environment variable definition; TODO: do something with it
			continue
		}
		// The first 6 tokens are min, hr, dom, mon, dow, user, and the
		// rest together form cmd
		// Convert and store the first 5 numbers (min, hr, dom, mon, dow)
		for i := 0; i < 5; i++ {
			err := m.convStoreToken(tokens[i], i)
			if err != nil {
				// Explicitly return this error
				return err
			}
		}
		// Store the 6th string (user)
		m.entries[m.entryCount].user = tokens[5]
		// Store the command
		m.entries[m.entryCount].cmd = strings.Join(tokens[6:], " ")
		// Increment the number of entries stored
		m.entryCount++
	}
	if err = scanner.Err(); err != nil {
		return
	}
	for i := 0; i < int(m.entryCount); i++ {
		fmt.Println(m.entries[i])
	}

	return
}

func (m *Manager) convStoreToken(token string, field int) (err error) {
	val, err := strconv.Atoi(token)
	if err != nil {
		return
	}

	// TODO: validation
	switch field {
	case 0:
		m.entries[m.entryCount].min = uint8(val)
	case 1:
		m.entries[m.entryCount].hr = uint8(val)
	case 2:
		m.entries[m.entryCount].dom = uint8(val)
	case 3:
		m.entries[m.entryCount].mon = uint8(val)
	case 4:
		m.entries[m.entryCount].dow = uint8(val)
	}

	return
}

// ServeHTTP serves HTTP requests
func (m *Manager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, m.crontabContents)
}
