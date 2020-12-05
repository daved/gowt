package toolsvc

import (
	"net/http"
)

//Index handler
func (s *ToolSvc) Index(w http.ResponseWriter, r *http.Request) {
	selDB, err := s.db.Query("SELECT * FROM tools ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	tool := toolData{}
	res := []toolData{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}
		s.log.Println("Listing Row: Id " + string(id) + " | name " + name + " | category " + category + " | url " + url + " | rating " + string(rating) + " | notes " + notes)

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
