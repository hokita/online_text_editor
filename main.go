package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
)

const (
	ExitCodeOk    = 0
	ExitCodeError = 1
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitCodeError)
	}

	os.Exit(ExitCodeOk)
}

func run() error {
	fmt.Println("start server")

	h := newHub()
	go h.run()

	file := newFile(filepath.Join("data", "text.txt"))

	r := mux.NewRouter()
	r.Handle("/", &initHandler{file: file})
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(h, file, w, r)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		return err
	}

	return nil
}

type initHandler struct {
	file *file
}

func (i *initHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(
		template.ParseFiles(filepath.Join("templates", "index.html")),
	)

	text, err := i.file.read()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	m := map[string]string{
		"Text": text,
		"Host": r.Host,
	}

	tpl.Execute(w, m)
}
