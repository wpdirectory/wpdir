package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/plugin"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/search"
	"github.com/wpdirectory/wpdir/internal/theme"
)

type errResponse struct {
	Code string `json:"code,omitempty"`
	Err  string `json:"error"`
}

// getSearches ...
func (s *Server) getSearches() http.HandlerFunc {
	type searchOverview struct {
		ID        string               `json:"id"`
		Input     string               `json:"input"`
		Repo      string               `json:"repo"`
		Matches   int                  `json:"matches"`
		Started   time.Time            `json:"started,omitempty"`
		Completed time.Time            `json:"completed,omitempty"`
		Progress  uint32               `json:"progress"`
		Status    search.Search_Status `json:"status"`
	}

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

// getSearch ...
func (s *Server) getSearch() http.HandlerFunc {
	type getSearchResponse struct {
		ID        string               `json:"id"`
		Input     string               `json:"input"`
		Repo      string               `json:"repo"`
		Matches   int                  `json:"matches"`
		Started   string               `json:"started,omitempty"`
		Completed string               `json:"completed,omitempty"`
		Progress  uint32               `json:"progress"`
		Status    search.Search_Status `json:"status"`
		Opts      search.Options       `json:"options"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if searchID := chi.URLParam(r, "id"); searchID != "" {
			if s.Manager.Exists(searchID) {
				srch := s.Manager.Get(searchID)
				writeResp(w, srch)
				return
			}

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

			writeResp(w, srch)
		} else {
			var resp errResponse
			resp.Err = "You must specify a valid Search ID."
			writeResp(w, resp)
		}
	}
}

// getSearchSummary ...
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

// getSearchMatches ...
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

// getMatchFile ...
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

// createSearch ...
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

// getRepo ...
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
				resp.CurrentRevision = s.Manager.Plugins.Rev()
			case "themes":
				resp.Name = repoName
				resp.Total = int(s.Manager.Themes.Len())
				resp.PendingUpdates = len(s.Manager.Themes.UpdateQueue)
				resp.CurrentRevision = s.Manager.Themes.Rev()
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
		Plugins *repo.Summary `json:"plugins,omitempty"`
		Themes  *repo.Summary `json:"themes,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp getRepoOverviewResponse

		resp.Plugins = s.Manager.Plugins.Summary()
		resp.Themes = s.Manager.Themes.Summary()

		writeResp(w, resp)
	}
}

// getPlugin ...
func (s *Server) getPlugin() http.HandlerFunc {
	type getPluginResponse struct {
		Slug                   string `json:"slug"`
		Name                   string `json:"name,omitempty"`
		Version                string `json:"version,omitempty"`
		Author                 string `json:"author,omitempty"`
		AuthorProfile          string `json:"author_profile,omitempty"`
		Rating                 int    `json:"rating,omitempty"`
		NumRatings             int    `json:"num_ratings,omitempty"`
		SupportThreads         int    `json:"support_threads,omitempty"`
		SupportThreadsResolved int    `json:"support_threads_resolved,omitempty"`
		ActiveInstalls         int    `json:"active_installs,omitempty"`
		Downloaded             int    `json:"downloaded,omitempty"`
		LastUpdated            string `json:"last_updated,omitempty"`
		Added                  string `json:"added,omitempty"`
		Homepage               string `json:"homepage,omitempty"`
		ShortDescription       string `json:"short_description,omitempty"`
		DownloadLink           string `json:"download_link,omitempty"`
		StableTag              string `json:"stable_tag,omitempty"`
		Status                 string `json:"status"`
		Err                    string `json:"error,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp getPluginResponse

		if slug := chi.URLParam(r, "slug"); slug != "" {
			p := s.Manager.Plugins.Get(slug).(*plugin.Plugin)
			resp.Slug = p.Slug
			resp.Name = p.Name
			resp.Version = p.Version
			resp.Author = p.Author
			resp.AuthorProfile = p.AuthorProfile
			resp.Rating = p.Rating
			resp.NumRatings = p.NumRatings
			resp.SupportThreads = p.SupportThreads
			resp.ActiveInstalls = p.ActiveInstalls
			resp.Downloaded = p.Downloaded
			resp.LastUpdated = p.LastUpdated
			resp.Added = p.Added
			resp.Homepage = p.Homepage
			resp.ShortDescription = p.ShortDescription
			resp.DownloadLink = p.DownloadLink
			resp.StableTag = p.StableTag
			resp.Status = p.GetStatus()
		} else {
			resp.Err = "You must specify a valid Plugin Name."
		}
		writeResp(w, resp)
	}
}

// getTheme ...
func (s *Server) getTheme() http.HandlerFunc {
	type getThemeResponse struct {
		Slug    string `json:"slug"`
		Name    string `json:"name,omitempty"`
		Version string `json:"version,omitempty"`
		Status  string `json:"status"`
		Err     string `json:"error,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp getThemeResponse

		if slug := chi.URLParam(r, "slug"); slug != "" {
			t := s.Manager.Themes.Get(slug).(*theme.Theme)
			resp.Slug = t.Slug
			resp.Name = t.Name
			resp.Version = t.Version
			resp.Status = t.GetStatus()
		} else {
			resp.Err = "You must specify a valid Theme Name."
		}

		writeResp(w, resp)
	}
}
