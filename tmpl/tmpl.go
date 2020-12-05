package tmpl

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
)

type Tmpl struct {
	*template.Template
	log *log.Logger
}

func New(log *log.Logger, t *template.Template) *Tmpl {
	return &Tmpl{
		Template: t,
		log:      log,
	}
}

func (t *Tmpl) WriteTemplate(w http.ResponseWriter, name string, v interface{}) {
	var b bytes.Buffer

	if err := t.Template.ExecuteTemplate(&b, name, v); err != nil {
		t.log.Printf("write template: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b.Bytes()); err != nil {
		t.log.Printf("write template: %v\n", err) // maybe notify ops if occurs
	}
}
