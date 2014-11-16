package manager

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

// ServeHTTP serves HTTP requests
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := "/index.htm"
	if r.URL.Path != "/" && r.URL.Path != page {
		page = r.URL.Path
	}
	// TODO: to see the pages, pwd should be github.com/teacoder/gron; fix this
	// TODO: parse body
	p := &Page{Title: "gron", Body: m.crontabContents}
	t, err := template.ParseFiles("manager/template" + page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func serveHTTPCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "manager/template"+r.URL.Path)
}
