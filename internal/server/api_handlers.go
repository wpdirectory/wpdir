package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/store/ulid"
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

	type getSearchResponse struct {
		Status  int                     `json:"status"`
		ID      string                  `json:"id"`
		Results []*index.SearchResponse `json:"results"`
		Err     string                  `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var resp getSearchResponse

		if searchID := chi.URLParam(r, "id"); searchID != "" {
			resp.ID = searchID
			s.lock.RLock()
			if val, ok := s.Searches[searchID]; ok {
				resp.Results = val
				resp.Status = 200
			} else {
				resp.Err = "Search Not Found."
				resp.Status = 404
			}
			s.lock.RUnlock()

		} else {
			resp.Err = "You must specify a valid Search ID."
			resp.Status = 404
		}

		writeResp(w, resp)
	}
}

// createSearch ...
func (s *Server) createSearch() http.HandlerFunc {
	type createSearchRequest struct {
		Input  string `json:"input"`
		Target string `json:"target"`
	}

	type createSearchResponse struct {
		Status int    `json:"status"`
		ID     string `json:"id"`
		Err    string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var data createSearchRequest
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}

		id := ulid.New()

		var empty []*index.SearchResponse

		s.lock.Lock()
		s.Searches[id] = empty
		s.lock.Unlock()

		// Perform non-blocking Search...
		go s.doSearch(data.Input, id)

		var resp createSearchResponse
		resp.ID = id
		writeResp(w, resp)
	}
}
