package manager

import (
	"fmt"
	"net/http"
)

// ServeHTTP serves HTTP requests
func (m *Manager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, m.crontabContents)
}
