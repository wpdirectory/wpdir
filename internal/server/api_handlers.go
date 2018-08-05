package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/search"
)

type errResponse struct {
	Code string `json:"code,omitempty"`
	Err  string `json:"error"`
}

// getSearches returns a list of between 10-100 public Searches
func (s *Server) getSearches() http.HandlerFunc {
	type getSearchesResponse struct {
		Searches []search.Search `json:"searches,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp getSearchesResponse

		searchLimit := chi.URLParam(r, "limit")
		limit, err := strconv.Atoi(searchLimit)
		if err != nil || limit < 10 || limit > 100 {
			var resp errResponse
			resp.Err = "You must specify a valid limit (10-100)."
			writeResp(w, resp)
			return
		}

		list := db.GetLatestPublicSearchList(limit)

		for _, id := range list {
			var srch search.Search
			bytes, err := db.GetSearch(id)
			if err != nil {
				continue
			}
			err = srch.Unmarshal(bytes)
			if err != nil {
				continue
			}
			resp.Searches = append(resp.Searches, srch)
		}

		writeResp(w, resp)
	}
}

// getSearch fetches the data for a search by ID
// First we check for inprogress searches in memory, then we check in the DB
func (s *Server) getSearch() http.HandlerFunc {
	type getSearchResponse struct {
		ID        string               `json:"id"`
		Input     string               `json:"input"`
		Repo      string               `json:"repo"`
		Matches   uint32               `json:"matches"`
		Started   string               `json:"started,omitempty"`
		Completed string               `json:"completed,omitempty"`
		Progress  uint32               `json:"progress"`
		Status    search.Search_Status `json:"status"`
		QueuePos  int                  `json:"queue_pos"`
		Opts      search.Options       `json:"options"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if searchID := chi.URLParam(r, "id"); searchID != "" {

			// Check InProgress Searches
			if s.Manager.Exists(searchID) {
				var resp getSearchResponse
				srch := s.Manager.Get(searchID)

				resp.ID = srch.ID
				resp.Input = srch.Input
				resp.Repo = srch.Repo
				resp.Matches = srch.Matches
				resp.Started = srch.Started
				resp.Completed = srch.Completed
				resp.Progress = srch.Progress
				resp.Status = srch.Status
				resp.QueuePos = s.Manager.Queue.Pos(searchID)
				resp.Opts = *srch.Options

				writeResp(w, resp)
				return
			}

			// Check Completed Searches in DB
			bytes, err := db.GetSearch(searchID)
			if err != nil || bytes == nil {
				var resp errResponse
				resp.Err = fmt.Sprintf("Search %s not found", searchID)
				writeResp(w, resp)
				return
			}

			var srch search.Search
			err = srch.Unmarshal(bytes)
			if err != nil {
				var resp errResponse
				resp.Err = fmt.Sprintf("Could not Unmarshal Search data for: %s\n", searchID)
				writeResp(w, resp)
				return
			}

			srch.Status = search.Completed

			writeResp(w, srch)
		} else {
			var resp errResponse
			resp.Err = "You must specify a valid Search ID."
			writeResp(w, resp)
		}
	}
}

// getSearchSummary returns a Summary of the Search results
func (s *Server) getSearchSummary() http.HandlerFunc {
	type getSearchSummaryResponse struct {
		Results []*search.Result `json:"results"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		searchID := chi.URLParam(r, "id")

		if searchID != "" {
			var resp getSearchSummaryResponse
			bytes, err := db.GetSummary(searchID)
			if err != nil || bytes == nil {
				var resp errResponse
				resp.Err = fmt.Sprintf("Summary not found for Search %s\n", searchID)
				writeResp(w, resp)
				return
			}

			var summary search.Summary
			err = summary.Unmarshal(bytes)
			if err != nil {
				var resp errResponse
				resp.Err = fmt.Sprintf("Could not Unmarshal Summary data for Search %s\n", searchID)
				writeResp(w, resp)
				return
			}

			for _, result := range summary.List {
				resp.Results = append(resp.Results, result)
			}

			writeResp(w, resp)
		} else {
			var resp errResponse
			resp.Err = "You must specify a valid Search ID."
			writeResp(w, resp)
		}
	}
}

// getSearchMatches returns a list of Search matches for a given Extension
func (s *Server) getSearchMatches() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchID := chi.URLParam(r, "id")
		slug := chi.URLParam(r, "slug")

		if searchID != "" && slug != "" {
			bytes, err := db.GetMatches(searchID, slug)
			if err != nil || bytes == nil {
				var resp errResponse
				resp.Err = fmt.Sprintf("Matches not found for Search %s and Slug %s\n", searchID, slug)
				writeResp(w, resp)
				return
			}

			var matches search.Matches
			err = matches.Unmarshal(bytes)
			if err != nil {
				var resp errResponse
				resp.Err = fmt.Sprintf("Could not Unmarshal Matches data for Search %s and Slug %s\n", searchID, slug)
				writeResp(w, resp)
				return
			}

			writeResp(w, matches)
		} else {
			var resp errResponse
			resp.Err = "You must specify a valid Search ID and Item Slug."
			writeResp(w, resp)
		}
	}
}

// getMatchFile returns the contents of a file identified by Repo, Slug and Filename
func (s *Server) getMatchFile() http.HandlerFunc {
	type getFileRequest struct {
		Repo string `json:"repo"`
		Slug string `json:"slug"`
		File string `json:"file"`
	}
	type getFileResponse struct {
		Code string `json:"code"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var data getFileRequest
		err := decoder.Decode(&data)
		if err != nil {
			var resp errResponse
			resp.Err = "Could not decode the POST body"
			writeResp(w, resp)
			return
		}

		if data.Repo != "" && data.Slug != "" && data.File != "" {
			var resp getFileResponse
			path, err := s.getFilePath(data.Repo, data.Slug, data.File)
			if err != nil {
				var resp errResponse
				resp.Err = "File could not be found"
				writeResp(w, resp)
				return
			}

			f, err := os.Open(path)
			if err != nil {
				var resp errResponse
				resp.Err = "File could not be opened"
				writeResp(w, resp)
				return
			}
			defer f.Close()

			// TODO: Consider making a sync.Pool of gzip Readers
			c, err := gzip.NewReader(f)
			if err != nil {
				var resp errResponse
				resp.Err = "File could not be decoded"
				writeResp(w, resp)
				return
			}
			defer c.Close()

			buf := new(bytes.Buffer)
			buf.ReadFrom(c)
			content := buf.String()

			resp.Code = string(content)
			writeResp(w, resp)
		} else {
			var resp errResponse
			resp.Err = "You must specify a valid repository, slug and filename"
			writeResp(w, resp)
		}
	}
}

// createSearch creates a new Search and returns the ID
func (s *Server) createSearch() http.HandlerFunc {
	type createSearchRequest struct {
		Input   string `json:"input"`
		Target  string `json:"target"`
		Private bool   `json:"private"`
	}

	type createSearchResponse struct {
		Status int    `json:"status"`
		ID     string `json:"id"`
		Err    string `json:"error,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp createSearchResponse
		decoder := json.NewDecoder(r.Body)

		var data createSearchRequest
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}

		// Ensure regex is not blank
		// TODO: Are there other checks which should be made here to prevent abuse?
		if data.Input == "" {
			var resp errResponse
			resp.Err = "Please provide non-blank search input."
			writeResp(w, resp)
			return
		}

		// Check Target
		switch data.Target {
		case "plugins":
			break
		case "themes":
			break
		default:
			var resp errResponse
			resp.Err = "Please provide a valid target"
			writeResp(w, resp)
			return
		}

		var sr search.Request
		sr.Input = data.Input
		sr.Repo = data.Target
		sr.Private = data.Private

		// Perform non-blocking Search...
		id := s.Manager.NewSearch(sr)

		resp.ID = id
		writeResp(w, resp)
	}
}

// getRepo returns an overview of the Repo identified by name
func (s *Server) getRepo() http.HandlerFunc {
	type getRepoResponse struct {
		Name            string `json:"name"`
		Total           int    `json:"total"`
		PendingUpdates  int    `json:"pending_updates"`
		CurrentRevision int    `json:"current_revision"`
		Err             string `json:"error,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var resp getRepoResponse

		if repoName := chi.URLParam(r, "name"); repoName != "" {
			switch repoName {
			case "plugins":
				resp.Name = repoName
				resp.Total = int(s.Manager.Plugins.Len())
				resp.PendingUpdates = len(s.Manager.Plugins.UpdateQueue)
				resp.CurrentRevision = s.Manager.Plugins.GetRev()
			case "themes":
				resp.Name = repoName
				resp.Total = int(s.Manager.Themes.Len())
				resp.PendingUpdates = len(s.Manager.Themes.UpdateQueue)
				resp.CurrentRevision = s.Manager.Themes.GetRev()
			default:
				resp.Err = "Repository Not Found."
			}
		} else {
			resp.Err = "You must specify a valid Repository Name."
		}

		writeResp(w, resp)
	}
}

// getRepoOverview ...
func (s *Server) getRepoOverview() http.HandlerFunc {
	type getRepoOverviewResponse struct {
		Plugins     *repo.Summary `json:"plugins,omitempty"`
		Themes      *repo.Summary `json:"themes,omitempty"`
		UpdateQueue int           `json:"update_queue,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp getRepoOverviewResponse

		resp.Plugins = s.Manager.Plugins.Summary()
		resp.Themes = s.Manager.Themes.Summary()
		resp.UpdateQueue = len(s.Manager.Plugins.UpdateQueue)

		writeResp(w, resp)
	}
}

// TODO: Combine the below into a single getExtension handler/endpoint

// getPlugin returns data for a Plugin Extension
func (s *Server) getPlugin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if slug := chi.URLParam(r, "slug"); slug != "" {
			p := s.Manager.Plugins.Get(slug)
			writeResp(w, p)
		} else {
			var resp errResponse
			resp.Err = "You must specify a valid Plugin Name"
			writeResp(w, resp)
		}
	}
}

// getTheme returns data for a Theme Extension
func (s *Server) getTheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if slug := chi.URLParam(r, "slug"); slug != "" {
			t := s.Manager.Themes.Get(slug)
			writeResp(w, t)
		} else {
			var resp errResponse
			resp.Err = "You must specify a valid Theme Name"
			writeResp(w, resp)
		}
	}
}
