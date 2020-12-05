package toolsvc

import (
	"net/http"
)

func (s *ToolSvc) Delete(w http.ResponseWriter, r *http.Request) {
	tool := r.URL.Query().Get("id")
	delForm, err := s.db.Prepare("DELETE FROM tools WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(tool)
	s.log.Println("DELETE " + tool)
	http.Redirect(w, r, "/", 301)
}
