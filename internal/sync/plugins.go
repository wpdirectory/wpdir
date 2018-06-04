package sync

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wpdirectory/wpdir/internal/index"
	"github.com/wpdirectory/wpdir/internal/repo"
	"github.com/wpdirectory/wpdir/internal/stats"
	"github.com/wpdirectory/wpdir/internal/store"
	"github.com/wpdirectory/wpdir/internal/svn"
)

func monitorPlugins(plugins chan svn.LogEntry) {

	var currentRevision = 0

	resp, err := repo.GetLogLatest("plugins", "")
	if err != nil {
		panic(err)
	}

	currentRevision = resp.Revision

	log.Printf("Beginning from revision %d\n", currentRevision)

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

func startPluginChecks(plugins chan svn.LogEntry) {

	// Start X number of workers reading from the plugins channel.
	for i := 1; i <= workerCount; i++ {

		go pluginWorker(plugins)

	}

}

func pluginWorker(plugins chan svn.LogEntry) {

	for {

		select {
		case plugin := <-plugins:

			// Extract the plugin slug
			parts := strings.Split(plugin.Paths[0].File, "/")
			slug := parts[1]

			actionable, path := repo.PluginNeedsUpdate(plugin)
			if actionable {

				log.Printf("%s (plugin) is being updated.\n", slug)

				// Update Plugin Files
				err := store.UpdatePlugin(path, slug)
				if err != nil {
					log.Printf("%s (plugin) could not be updated: %s\n", slug, err)
				}

				log.Printf("%s (plugin) is being indexed.\n", slug)

				// Update Plugin Index
				opts := &index.IndexOptions{
					ExcludeDotFiles: true,
				}
				wd, _ := os.Getwd()
				src := filepath.Join(wd, "data", "plugins", slug)
				dest := filepath.Join(wd, "index", "plugins", slug)
				url := slug
				rev := plugin.Revision
				_, err = index.BuildNew(opts, dest, src, url, string(rev))
				if err != nil {
					log.Printf("%s (plugin) could not be indexed: %s\n", slug, err)
				}

			} else {

				log.Printf("%s (plugin) does not require action.\n", slug)

			}

			stats.AddLatestPlugin(slug, plugin.Revision, plugin.Date)

		}

	}

}
