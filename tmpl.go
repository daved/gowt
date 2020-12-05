package main

import (
	"bytes"
	"net/http"
	"text/template"
)

type tmpl struct {
	*template.Template
}

func (t *tmpl) WriteTemplate(w http.ResponseWriter, name string, v interface{}) {
	var b bytes.Buffer

	if err := t.Template.ExecuteTemplate(&b, name, v); err != nil {
		log.Printf("write template: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Printf("write template: %v\n", err) // maybe notify ops if occurs
	}
}
