package toolsvc

import "net/http"

func (s *ToolSvc) New(w http.ResponseWriter, r *http.Request) {
	s.t.WriteTemplate(w, "New", nil)
}
