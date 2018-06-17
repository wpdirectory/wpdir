package repo

import (
	"fmt"
	"time"

	"github.com/wpdirectory/wpdir/internal/svn"
)

const (
	// WPRepoURL is the URL for the WordPress SVN Repos
	WPRepoURL = "https://%s.svn.wordpress.org/%s"
)

var (
	pluginManagementUser = "plugin-master"
	themeManagementUser  = "theme-master"
)

// Extension represents a WordPress Extension (Plugin/Theme).
type Extension struct {
	Name     string
	Slug     string
	Version  string
	Revision string
	Dir      string
	Time     time.Time
}

// GetLogLatest gets details of the latest Revision.
func GetLogLatest(repo string, path string) (svn.LogEntry, error) {

	URL := fmt.Sprintf(WPRepoURL, repo, path)

	args := []string{URL, "-v", "-r", "HEAD"}

	out, err := svn.Log(args...)

	if len(out) > 0 {
		return out[0], err
	}

	return svn.LogEntry{}, err

}

// GetLogDiff gets details of the Revisions between the two values provided.
func GetLogDiff(repo string, path string, start int, end int) ([]svn.LogEntry, error) {

	URL := fmt.Sprintf(WPRepoURL, repo, path)

	diff := fmt.Sprintf("%d:%d", start, end)

	args := []string{URL, "-v", "-r", diff}

	out, err := svn.Log(args...)

	return out, err

}

// GetList gets details of all folders inside the path.
func GetList(repo string, path string) ([]svn.ListEntry, error) {

	URL := fmt.Sprintf(WPRepoURL, repo, path)

	args := []string{URL}

	out, err := svn.List(args...)

	return out, err

}

// GetDiff gets the files which changed between the two values provided.
func GetDiff(repo string, path string, start int, end int) ([]byte, error) {

	URL := fmt.Sprintf(WPRepoURL, repo, path)
	diff := fmt.Sprintf("%d:%d", start, end)
	args := []string{URL, "-r", diff, "--xml", "--summarize"}

	out, err := svn.Diff(args...)

	return out, err

}

// GetCat gets the contents from the remote path, at (optional) revision.
func GetCat(repo string, path string, revision int) ([]byte, error) {

	URL := fmt.Sprintf(WPRepoURL, repo, path)
	args := []string{URL}

	if revision != 0 {
		args = append(args, "-r")
		args = append(args, fmt.Sprintf("%d", revision))
	}

	out, err := svn.Cat(args...)

	return out, err

}

// DoExport writes the remote files from path to local dest.
func DoExport(repo string, path string, dest string) error {

	URL := fmt.Sprintf(WPRepoURL, repo, path)
	args := []string{URL, dest, "-q", "--force", "--depth", "infinity"}

	err := svn.Export(args...)

	return err

}

// DoCheckout creates a local svn repo at dest from remote at path.
func DoCheckout(repo string, path string, dest string) error {

	URL := fmt.Sprintf(WPRepoURL, repo, path)
	args := []string{URL, dest, "-q"}

	err := svn.Checkout(args...)

	return err

}
