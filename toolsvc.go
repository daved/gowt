package main

import (
	"database/sql"
	"net/http"
	"text/template"
)

// Tool struct
type Tool struct {
	Id       int
	Name     string
	Category string
	URL      string
	Rating   int
	Notes    string
}

type toolSvc struct {
	db *sql.DB
	t  *tmpl
}

func newToolSvc(db *sql.DB, t *template.Template) *toolSvc {
	return &toolSvc{
		db: db,
		t:  &tmpl{t},
	}
}

//Index handler
func (s *toolSvc) Index(w http.ResponseWriter, r *http.Request) {
	selDB, err := s.db.Query("SELECT * FROM tools ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	tool := Tool{}
	res := []Tool{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}
		log.Println("Listing Row: Id " + string(id) + " | name " + name + " | category " + category + " | url " + url + " | rating " + string(rating) + " | notes " + notes)

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
		res = append(res, tool)
	}

	s.t.WriteTemplate(w, "Index", res)
}

//Show handler
func (s *toolSvc) Show(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")
	selDB, err := s.db.Query("SELECT * FROM tools WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	defer selDB.Close() //nolint

	tool := Tool{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}

		log.Println("Showing Row: Id " + string(id) + " | name " + name + " | category " + category + " | url " + url + " | rating " + string(rating) + " | notes " + notes)

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
	}

	s.t.WriteTemplate(w, "Show", tool)
}

func (s *toolSvc) New(w http.ResponseWriter, r *http.Request) {
	s.t.WriteTemplate(w, "New", nil)
}

func (s *toolSvc) Edit(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")
	selDB, err := s.db.Query("SELECT * FROM tools WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	tool := Tool{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
	}

	s.t.WriteTemplate(w, "Edit", tool)
}

func (s *toolSvc) Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		category := r.FormValue("category")
		url := r.FormValue("url")
		rating := r.FormValue("rating")
		notes := r.FormValue("notes")
		insForm, err := s.db.Prepare("INSERT INTO tools (name, category, url, rating, notes) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, category, url, rating, notes)
		log.Println("Insert Data: name " + name + " | category " + category + " | url " + url + " | rating " + rating + " | notes " + notes)
	}
	http.Redirect(w, r, "/", 301)
}

func (s *toolSvc) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		category := r.FormValue("category")
		url := r.FormValue("url")
		rating := r.FormValue("rating")
		notes := r.FormValue("notes")
		id := r.FormValue("uid")
		insForm, err := s.db.Prepare("UPDATE tools SET name=?, category=?, url=?, rating=?, notes=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, category, url, rating, notes, id)
		log.Println("UPDATE Data: name " + name + " | category " + category + " | url " + url + " | rating " + rating + " | notes " + notes)
	}
	http.Redirect(w, r, "/", 301)
}

func (s *toolSvc) Delete(w http.ResponseWriter, r *http.Request) {
	tool := r.URL.Query().Get("id")
	delForm, err := s.db.Prepare("DELETE FROM tools WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(tool)
	log.Println("DELETE " + tool)
	http.Redirect(w, r, "/", 301)
}
