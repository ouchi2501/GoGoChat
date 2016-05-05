package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// templ is once template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// serveHTTP is job to HTTP Request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "application addr")
	flag.Parse()
	r := newRoom()
	// root
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// started chatRoom
	go r.run()
	// started webServer
	log.Println("started web server... port", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
