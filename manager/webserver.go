package manager

import (
	"html/template"
	"net/http"
)

// ServeHTTP serves HTTP requests
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := "/index.htm"
	if r.URL.Path != "/" && r.URL.Path != page {
		page = r.URL.Path
	}
	// TODO: to see the pages, pwd should be github.com/1d4Nf6/zinc; fix this
	t, err := template.ParseFiles("manager/template" + page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func serveHTTPCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "manager/template"+r.URL.Path)
}
