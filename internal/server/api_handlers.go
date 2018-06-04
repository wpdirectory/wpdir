package server

import "net/http"

func (s *Server) getSearchList() http.HandlerFunc {
	type fightsResponse struct {
		Fights []string
		Time   string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		fights := fightsResponse{
			Fights: []string{
				"Fight 1",
				"Fight 2",
				"Fight 3",
				"Fight 4",
				"Fight 5",
				"Fight 6",
				"Fight 7",
				"Fight 8",
			},
			Time: "15:59 12/05/2018",
		}
		writeResp(w, fights)
	}
}

func (s *Server) getSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fight := ""
		writeResp(w, fight)
	}
}
func (s *Server) createSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fight := ""
		writeResp(w, fight)
	}
}
