package toolsvc

import (
	"net/http"
)

func (s *ToolSvc) Insert(w http.ResponseWriter, r *http.Request) {
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
		s.log.Println("Insert Data: name " + name + " | category " + category + " | url " + url + " | rating " + rating + " | notes " + notes)
	}
	http.Redirect(w, r, "/", 301)
}

func (s *ToolSvc) Update(w http.ResponseWriter, r *http.Request) {
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
		s.log.Println("UPDATE Data: name " + name + " | category " + category + " | url " + url + " | rating " + rating + " | notes " + notes)
	}
	http.Redirect(w, r, "/", 301)
}
