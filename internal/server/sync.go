package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/search"
	"github.com/wpdirectory/wpdir/internal/store"
	"github.com/wpdirectory/wpdir/internal/store/ulid"
	"github.com/wpdirectory/wpdir/internal/svn"
)

var (
	storageDir  string
	indexDir    string
	workerCount int
	bufferSize  int
)

func (s *Server) startSync() {

	// Check if SVN client is available
	if !svn.IsClientInstalled() {
		log.Fatal("The SVN client is not available, it is required for operation.")
	}

	// Setup Config Based Values
	storageDir = "data"
	indexDir = "index"
	workerCount = 10
	bufferSize = 90000

	// Create Plugins Channel
	plugins := make(chan svn.LogEntry, bufferSize)

	// Load existing indexes
	s.loadPluginIndexes()

	// Test Search
	//s.doSearch("pre_get_posts")

	// Fresh start, queue all plugins for indexing.
	//s.freshSlurp()
	//s.indexSlurp()
	//fmt.Println("Finished Indexing")
	//os.Exit(1)
	//s.freshLoadPlugins(plugins)

	// Monitor plugins directory for updates and spawn workers
	// to process updates.
	go s.pluginsMonitor(plugins)
	go s.pluginWorkerStartup(plugins)

}

func (s *Server) doSearch(input string, id string) {
	var filesOpened int
	var durationMs int
	var results []*index.SearchResponse

	results, err := searchAll(input, s.PluginSearchers, &filesOpened, &durationMs)
	if err != nil {
		// TODO(knorton): Return ok status because the UI expects it for now.
		s.Logger.Println("Search Failed: ", err)
		return
	}
	s.Logger.Println("Search Finished")

	var matches int
	for _, result := range results {
		for _, match := range result.Matches {
			matches += len(match.Matches)
		}
	}

	fmt.Printf("Total Plugins: %d\n", len(results))
	fmt.Printf("Files Opened: %d\n", filesOpened)
	fmt.Printf("Matches: %d\n", matches)
	fmt.Printf("Duration: %d\n", durationMs)

	s.lock.RLock()
	s.Searches[id] = results
	s.lock.RUnlock()

}

func searchAll(query string, idx map[string]*search.Searcher, filesOpened *int, durationMs *int) ([]*index.SearchResponse, error) {

	startedAt := time.Now()

	n := len(idx)

	opts := &index.SearchOptions{
		Offset:         0,
		Limit:          0,
		LinesOfContext: 0,
		IgnoreCase:     false,
	}

	limiter := make(chan struct{}, runtime.NumCPU())

	// use a buffered channel to avoid routine leaks on errs.
	ch := make(chan *index.SearchResponse, n)
	for slug := range idx {
		limiter <- struct{}{}
		go func(slug string) {
			fms, err := idx[slug].Search(query, opts)

			if err != nil {
				fmt.Println("Search Error: ", err)
				<-limiter
				return
			}
			fms.Slug = slug
			<-limiter
			ch <- fms
			return
		}(slug)
	}

	res := []*index.SearchResponse{}
	for i := 0; i < n; i++ {
		r := <-ch
		if len(r.Matches) > 0 {
			res = append(res, r)
		}
		*filesOpened += r.FilesOpened
	}

	*durationMs = int(time.Now().Sub(startedAt).Seconds() * 1000)

	return res, nil
}

func (s *Server) pluginsMonitor(plugins chan svn.LogEntry) {

	// Get Current Revision
	var currentRevision = 0

	resp, err := repo.GetLogLatest("plugins", "")
	if err != nil {
		s.Logger.Fatalf("Failed to get latest revision from the Plugins repo: %s\n", err)
	}
	currentRevision = resp.Revision

	s.Logger.Printf("Beginning from revision %d\n", currentRevision)

	for {
		// Sleep 15 seconds so we do not hit the WP.org SVN too hard.
		time.Sleep(15 * time.Second)

		// Request the latest revision.
		resp, err := repo.GetLogLatest("plugins", "")
		if err != nil {
			// We need this value, restart loop.
			continue
		}

		// If the received revision is greater than current updates have been made.
		if resp.Revision > currentRevision {

			// TODO: Break this into multiple calls if the difference is large enough.
			// Prevent hitting the SVN repo with too large a request, which is currently 100 as
			// shown at: https://meta.trac.wordpress.org/browser/sites/trunk/wordpress.org/public_html/wp-content/plugins/plugin-directory/cli/class-svn-watcher.php#L30

			// Fetch the changelog between the current and latest revision.
			// Add one to currentRevision to avoid repeated changes.
			changelog, err := repo.GetLogDiff("plugins", "", currentRevision+1, resp.Revision)
			if err != nil {
				// We need these values, restart loop.
				continue
			}

			var counter int

			// Add each separate plugin which has been updated to the channel.
			for _, plugin := range changelog {
				plugins <- plugin
				counter++
			}

			// We have updated to the new revision, set as current.
			currentRevision = resp.Revision

		}

	}

}

func (s *Server) freshLoadPlugins(plugins chan svn.LogEntry) {

	resp, err := repo.GetLogLatest("plugins", "")
	if err != nil {
		s.Logger.Fatalf("Failed to get latest revision from the Plugins repo: %s\n", err)
	}

	rev := string(resp.Revision)

	list, err := repo.GetList("plugins", "")
	if err != nil {
		s.Logger.Fatalf("Could not fetch full list of Plugins: %s\n", err)
		return
	}

	limiter := make(chan struct{}, 32)

	for _, item := range list {

		limiter <- struct{}{}
		go func(item svn.ListEntry) {

			slug := item.Name
			id := ulid.New()
			tmpDir := filepath.Join(s.Config.WD, "data", "tmp", "plugins", slug)
			destDir := filepath.Join(s.Config.WD, "data", "index", "plugins", id)

			// Get new files and store in tmp dir
			err := store.GetPlugin(slug+"/trunk", slug)
			if err != nil {
				s.Logger.Printf("%s (plugin) could not be updated: %s\n", slug, err)
				store.DeleteFolder(tmpDir)
				<-limiter
				return
			}

			// Walk dir to get file stats
			// Num Files / Filesize / Largest Files / etc

			// Get ID and Generate Index using the ID for folder name
			opts := &index.IndexOptions{
				ExcludeDotFiles: true,
			}
			ref, err := index.Build(opts, destDir, tmpDir, slug, rev)
			if err != nil {
				s.Logger.Printf("%s (plugin) could not be indexed: %s\n", slug, err)
				store.DeleteFolder(tmpDir)
				store.DeleteFolder(destDir)
				<-limiter
				return
			}

			srchr, err := search.New(slug, slug, rev, ref)
			if err != nil {
				s.Logger.Printf("Could not create Searcher for %s (plugin): %s\n", slug, err)
				store.DeleteFolder(tmpDir)
				store.DeleteFolder(destDir)
				<-limiter
				return
			}

			s.lock.RLock()
			s.PluginSearchers[slug] = srchr
			s.lock.RUnlock()

			// Success: lock index and switch to new (symlink dir and new index reference) and add stats to DB
			// Failure: delete tmp/index files and log error
			store.DeleteFolder(tmpDir)

			<-limiter
		}(item)

	}

}

func (s *Server) pluginWorkerStartup(plugins chan svn.LogEntry) {

	for {

		select {
		case plugin := <-plugins:

			id := ulid.New()
			err := s.processPluginUpdate(plugin, id)
			if err != nil {
				s.Logger.Printf("Update Failed (plugin: %s): %s\n", plugin.Paths[0].File, err)
			}

		}

	}

	//for i := 1; i <= workerCount; i++ {

	//go s.pluginsUpdater(plugins)

	//}

}

func (s *Server) pluginsUpdater(plugins chan svn.LogEntry) {

	for {

		select {
		case plugin := <-plugins:

			id := ulid.New()
			err := s.processPluginUpdate(plugin, id)
			if err != nil {
				s.Logger.Printf("Update Failed (plugin: %s): %s\n", plugin.Paths[0].File, err)
			}

		}

	}

}

// UpdatePlugin ...
func (s *Server) processPluginUpdate(plugin svn.LogEntry, id string) error {

	// Prepare Information about Plugin
	parts := strings.Split(plugin.Paths[0].File, "/")
	slug := parts[1]
	tmpDir := filepath.Join(s.Config.WD, "data", "tmp", "plugins", slug)
	destDir := filepath.Join(s.Config.WD, "data", "index", "plugins", id)
	rev := string(plugin.Revision)

	// Check if we need to update local files
	needs, path := repo.PluginNeedsUpdate(plugin)
	if !needs {
		return nil
	}

	// Get new files and store in tmp dir
	err := store.GetPlugin(path, slug)
	if err != nil {
		s.Logger.Printf("%s (plugin) could not be updated: %s\n", slug, err)
		store.DeleteFolder(tmpDir)
		return err
	}

	// Walk dir to get file stats
	// Num Files / Filesize / Largest Files / etc

	// Get ID and Generate Index using the ID for folder name
	opts := &index.IndexOptions{
		ExcludeDotFiles: true,
	}
	ref, err := index.Build(opts, destDir, tmpDir, slug, string(plugin.Revision))
	if err != nil {
		s.Logger.Printf("%s (plugin) could not be indexed: %s\n", slug, err)
		store.DeleteFolder(tmpDir)
		store.DeleteFolder(destDir)
		return err
	}

	_, ok := s.PluginSearchers[slug]
	if ok {

		idx, err := ref.Open()
		if err != nil {
			s.Logger.Printf("Could not build new index for %s (plugin): %s\n", slug, err)
			store.DeleteFolder(tmpDir)
			store.DeleteFolder(destDir)
			return err
		}
		s.PluginSearchers[slug].SwapIndexes(idx)

	} else {

		srchr, err := search.New(slug, slug, rev, ref)
		if err != nil {
			s.Logger.Printf("Could not create Searcher for %s (plugin): %s\n", slug, err)
			store.DeleteFolder(tmpDir)
			store.DeleteFolder(destDir)
			return err
		}

		s.lock.RLock()
		s.PluginSearchers[slug] = srchr
		s.lock.RUnlock()

	}

	// Success: lock index and switch to new (symlink dir and new index reference) and add stats to DB
	// Failure: delete tmp/index files and log error
	store.DeleteFolder(tmpDir)

	return nil

}

func (s *Server) loadPluginIndexes() {

	indexDir := filepath.Join(s.Config.WD, "data", "index", "plugins")

	dirs, err := ioutil.ReadDir(indexDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range dirs {
		// If not Directory discard.
		if !dir.IsDir() {
			continue
		}

		// Read index
		ref, err := index.Read(filepath.Join(indexDir, dir.Name()))
		if err != nil {
			s.Logger.Printf("Could not read index directory: %s\n", err)
		}

		// Create Searcher
		srchr, err := search.New(ref.Slug, ref.Slug, ref.Rev, ref)
		if err != nil {
			continue
		}

		s.lock.RLock()
		s.PluginSearchers[ref.Slug] = srchr
		s.lock.RUnlock()

	}

}
