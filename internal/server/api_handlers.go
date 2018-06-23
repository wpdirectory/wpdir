package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/wpdirectory/wpdir/internal/plugin"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/theme"
)

type errResponse struct {
	Code string `json:"code,omitempty"`
	Err  string `json:"error"`
}

// getSearches ...
func (s *Server) getSearches() http.HandlerFunc {
	type searchOverview struct {
		ID        string    `json:"id"`
		Input     string    `json:"input"`
		Repo      string    `json:"repo"`
		Matches   int       `json:"matches"`
		Started   time.Time `json:"started,omitempty"`
		Completed time.Time `json:"completed,omitempty"`
		Progress  int       `json:"progress"`
		Total     int       `json:"total"`
		Status    status    `json:"status"`
	}

	type getSearchesResponse struct {
		Searches []*searchOverview `json:"searches,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp getSearchesResponse
		s.Searches.Lock()
		defer s.Searches.Unlock()

		searchLimit := chi.URLParam(r, "limit")
		limit, err := strconv.Atoi(searchLimit)
		if err != nil || limit < 10 || limit > 100 {
			var resp errResponse
			resp.Err = "You must specify a valid limit (10-100)."
			writeResp(w, resp)
			return
		}

		i := 1
		for _, srch := range s.Searches.List {
			so := &searchOverview{
				ID:        srch.ID,
				Input:     srch.Input,
				Repo:      srch.Repo,
				Matches:   len(srch.Matches),
				Started:   srch.Started,
				Completed: srch.Completed,
				Progress:  srch.Progress,
				Total:     srch.Total,
				Status:    srch.Status,
			}
			resp.Searches = append(resp.Searches, so)
			if i++; i == limit {
				break
			}
		}

		writeResp(w, resp)
	}
}

// getSearch ...
func (s *Server) getSearch() http.HandlerFunc {
	type getSearchResponse struct {
		ID        string    `json:"id"`
		Input     string    `json:"input"`
		Repo      string    `json:"repo"`
		Matches   int       `json:"matches"`
		Started   time.Time `json:"started,omitempty"`
		Completed time.Time `json:"completed,omitempty"`
		Progress  int       `json:"progress"`
		Total     int       `json:"total"`
		Status    status    `json:"status"`
		Opts      Options   `json:"options"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if searchID := chi.URLParam(r, "id"); searchID != "" {

			var resp getSearchResponse
			s.Searches.Lock()
			defer s.Searches.Unlock()
			srch, ok := s.Searches.List[searchID]
			if !ok {
				var resp errResponse
				resp.Err = fmt.Sprintf("Search %s not found.", searchID)
				writeResp(w, resp)
				return
			}
			srch.Lock()
			defer srch.Unlock()
			resp.ID = srch.ID
			resp.Input = srch.Input
			resp.Repo = srch.Repo
			resp.Matches = len(srch.Matches)
			if !srch.Started.IsZero() {
				resp.Started = srch.Started
			}
			if !srch.Completed.IsZero() {
				resp.Completed = srch.Completed
			}
			resp.Progress = srch.Progress
			resp.Total = srch.Total
			resp.Status = srch.Status
			resp.Opts = srch.Opts

			writeResp(w, resp)

		} else {
			var resp errResponse
			resp.Err = "You must specify a valid Search ID."
			writeResp(w, resp)
		}
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

		var sr SearchRequest
		sr.Input = data.Input
		sr.Repo = data.Target

		s.Searches.Lock()
		//s.Searches[id] = empty
		s.Searches.Unlock()

		// Perform non-blocking Search...
		id := s.Searches.NewSearch(sr)

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
				resp.Total = s.Plugins.Len()
				resp.PendingUpdates = len(s.Plugins.(*repo.PluginRepo).UpdateQueue)
				resp.CurrentRevision = s.Plugins.Rev()
			case "themes":
				resp.Name = repoName
				resp.Total = s.Themes.Len()
				resp.PendingUpdates = len(s.Themes.(*repo.ThemeRepo).UpdateQueue)
				resp.CurrentRevision = s.Themes.Rev()
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
		Plugins *repo.RepoSummary `json:"plugins,omitempty"`
		Themes  *repo.RepoSummary `json:"themes,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp getRepoOverviewResponse

		resp.Plugins = s.Plugins.Summary()
		resp.Themes = s.Themes.Summary()

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
			p := s.Plugins.Get(slug).(*plugin.Plugin)
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
			t := s.Themes.Get(slug).(*theme.Theme)
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
