package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	xlog "log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"

	"github.com/daved/gowt/dbase"
	"github.com/daved/gowt/toolsvc"
)

var log = xlog.New(ioutil.Discard, "", xlog.LstdFlags)

func main() {
	var (
		debug bool
	)

	flag.BoolVar(&debug, "debug", debug, "turn on debug logging")
	flag.Parse()

	if debug {
		log = xlog.New(os.Stdout, "", xlog.LstdFlags)
	}

	log.Println("Server started on: http://localhost:8080")

	db, err := dbase.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	t := template.New("all")
	t, err = t.ParseGlob("templates/*")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	tSvc := toolsvc.New(log, db, t)

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(tSvc.Index))
	mux.Handle("/show", http.HandlerFunc(tSvc.Show))
	mux.Handle("/new", http.HandlerFunc(tSvc.New))
	mux.Handle("/edit", http.HandlerFunc(tSvc.Edit))
	mux.Handle("/insert", http.HandlerFunc(tSvc.Insert))
	mux.Handle("/update", http.HandlerFunc(tSvc.Update))
	mux.Handle("/delete", http.HandlerFunc(tSvc.Delete))

	http.ListenAndServe(":8080", mux)
}
