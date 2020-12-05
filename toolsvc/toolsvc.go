package toolsvc

import (
	"log"
	"text/template"

	"github.com/daved/gowt/dbase"
	"github.com/daved/gowt/tmpl"
)

type ToolSvc struct {
	log *log.Logger
	db  *dbase.DBase
	t   *tmpl.Tmpl
}

func New(log *log.Logger, db *dbase.DBase, t *template.Template) *ToolSvc {
	return &ToolSvc{
		log: log,
		db:  db,
		t:   tmpl.New(log, t),
	}
}

// toolData struct
type toolData struct {
	Id       int
	Name     string
	Category string
	URL      string
	Rating   int
	Notes    string
}
