package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/wpdirectory/wpdir/internal/db"
	"github.com/wpdirectory/wpdir/internal/plugin"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/theme"
)

type errResponse struct {
	Code string `json:"code,omitempty"`
	Err  string `json:"error"`
}

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

// getSearchesLatest ...
func (s *Server) getSearchesLatest() http.HandlerFunc {
	type getSearchesLatestResponse struct {
		Searches []*LatestSearch `json:"searches,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var resp getSearchesLatestResponse

		resp.Searches = s.Searches.Latest.Get()

		writeResp(w, resp)
	}
}

// getSearch ...
func (s *Server) getSearch() http.HandlerFunc {
	type getSearchResponse struct {
		ID        string    `json:"id"`
		Input     string    `json:"input"`
		Repo      string    `json:"repo"`
		Matches   []*Match  `json:"matches"`
		Started   time.Time `json:"started"`
		Completed time.Time `json:"completed,omitempty"`
		Progress  int       `json:"progress"`
		Total     int       `json:"total"`
		Status    status    `json:"status"`
		Opts      Options   `json:"options"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if searchID := chi.URLParam(r, "id"); searchID != "" {

			bytes, err := db.GetFromBucket(searchID, "searches")
			if err == nil {
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				w.Header().Set("Vary", "Accept-Encoding")
				w.WriteHeader(http.StatusOK)
				w.Write(bytes)
				return
			}

			var resp getSearchResponse
			s.Searches.Lock()
			defer s.Searches.Unlock()
			srch := s.Searches.List[searchID]
			srch.Lock()
			defer srch.Unlock()
			resp.ID = srch.ID
			resp.Input = srch.Input
			resp.Repo = srch.Repo
			resp.Matches = srch.Matches
			resp.Started = srch.Started
			resp.Completed = srch.Completed
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
