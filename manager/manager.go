/*
Package manager implements the gron manager.
*/
package manager

import (
	"fmt"
)

// Manager defines the manager struct
type Manager struct {
	crontabPath string
	portNum     uint16
}

// CrontabPath is the getter function for crontabPath
func (m *Manager) CrontabPath() string {
	return m.crontabPath
}

// PortNum returns the port number
func (m *Manager) PortNum() uint16 {
	return m.portNum
}

// NewManager Creates a new manager
func NewManager() *Manager {
	return &Manager{crontabPath: "gron.go", portNum: 40000}
}

// ReadCrontab reads the crontab
func ReadCrontab(crontabPath string) (x int, err error) {
	x, err = 0, nil
	m := NewManager()

	fmt.Printf("Reading crontab file %s...\n", m.CrontabPath())
	fmt.Println(m.PortNum())

	return
}
