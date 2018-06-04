package server

import (
	"encoding/json"
	"net/http"
)

// getSearchList ...
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

// getSearch ...
func (s *Server) getSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fight := ""
		writeResp(w, fight)
	}
}

// createSearch ...
func (s *Server) createSearch() http.HandlerFunc {
	type createSearchRequest struct {
		input  string
		target string
	}

	type createSearchResponse struct {
		status int
		id     string
		err    string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var data createSearchRequest
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}

		// Create Search...

		var resp createSearchResponse
		writeResp(w, resp)
	}
}
