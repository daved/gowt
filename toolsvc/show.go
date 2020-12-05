package toolsvc

import (
	"net/http"
)

//Show handler
func (s *ToolSvc) Show(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")
	selDB, err := s.db.Query("SELECT * FROM tools WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	defer selDB.Close() //nolint

	tool := toolData{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}

		s.log.Println("Showing Row: Id " + string(id) + " | name " + name + " | category " + category + " | url " + url + " | rating " + string(rating) + " | notes " + notes)

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
	}

	s.t.WriteTemplate(w, "Show", tool)
}
