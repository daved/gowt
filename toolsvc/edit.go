package toolsvc

import "net/http"

func (s *ToolSvc) Edit(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")
	selDB, err := s.db.Query("SELECT * FROM tools WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	tool := toolData{}

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
