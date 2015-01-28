/*
Package manager implements the zinc manager.
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

// @jayanthc
// The Manager struct variables are private to the package. We need a getter.
func (e Entry) Get() map[string]interface{} {
	m := map[string]interface{}{
		"min": e.min, "hr": e.hr, "dom": e.dom, "mon": e.mon,
		"dow": e.dow, "user": e.user, "cmd": e.cmd}
	return m
}

// Manager defines the manager struct
type Manager struct {
	// Crontab-related variables
	crontab         string
	crontabContents []byte
	entries         [MaxEntries]Entry
	entryCount      uint8 // TODO: think of a better data type
	// Webserver-related variables
	hostname string
	port     uint16
}

// New Creates a new manager and performs initialization
func New(userCrontab string) (m *Manager, err error) {
	crontab := defCrontab
	if userCrontab != "" {
		crontab = userCrontab
	}
	// Create and initialize manager
	m = &Manager{crontab: crontab, hostname: defHostname, port: defPort}
	err = m.readCrontab()
	if err != nil {
		return m, err
	}

	// Run webserver in a goroutine
	http.Handle("/", m)
	http.HandleFunc("/css/", serveHTTPCSS)
	go http.ListenAndServe(m.hostname+":"+strconv.Itoa(int(m.port)), nil)

	return m, err
}

// Crontab gets the crontab file name, for the template
func (m *Manager) Crontab() string {
	return m.crontab
}

// Entries gets the crontab entries stored in the Manager
func (m *Manager) Entries() [MaxEntries]Entry {
	return m.entries
}

// EntryCount gets the number of crontab entries stored in the Manager
func (m *Manager) EntryCount() uint8 {
	return m.entryCount
}

//BuildEntry builds a crontab entry from the Entry struct
// TODO: check why this should be 'en Entry' and not 'en *Entry'
func (en Entry) BuildEntry() string {
	return fmt.Sprintf("%d\t%d\t%d\t%d\t%d\t%s\t%s",
		en.min, en.hr, en.dom, en.mon, en.dow, en.user, en.cmd)
}

// readCrontab reads the crontab and builds the entries data structure
func (m *Manager) readCrontab() (err error) {
	err = nil

	// Read the crontab
	fmt.Printf("Reading crontab file %s...\n", m.crontab)
	m.crontabContents, err = ioutil.ReadFile(m.crontab)
	if err != nil {
		return
	}

	// Scan the crontab lines
	// TODO: is there a simpler way to do this?
	scanner := bufio.NewScanner(strings.NewReader(string(m.crontabContents)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			// Empty line; Ignore
			continue
		}
		// Replace wildcard with "255"
		line = strings.Replace(line, "*", "255", -1)
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
