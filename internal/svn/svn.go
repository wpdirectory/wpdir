package svn

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	// WPRepoURL is the URL for the WordPress SVN Repos
	WPRepoURL = "https://%s.svn.wordpress.org/%s"
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
func GetLogLatest(repo string, path string) (LogEntry, error) {

	URL := fmt.Sprintf(WPRepoURL, repo, path)

	args := []string{URL, "-v", "-r", "HEAD"}

	out, err := log(args...)

	if len(out) > 0 {
		return out[0], err
	}

	return LogEntry{}, err

}

// GetLogDiff gets details of the Revisions between the two values provided.
func GetLogDiff(repo string, path string, start int, end int) ([]LogEntry, error) {

	URL := fmt.Sprintf(WPRepoURL, repo, path)

	diff := fmt.Sprintf("%d:%d", start, end)

	args := []string{URL, "-v", "-r", diff}

	out, err := log(args...)

	return out, err

}

// GetList gets details of all folders inside the path.
func GetList(repo string, path string) ([]ListEntry, error) {

	URL := fmt.Sprintf(WPRepoURL, repo, path)

	args := []string{URL}

	out, err := list(args...)

	return out, err

}

// GetDiff gets the files which changed between the two values provided.
func GetDiff(repo string, path string, start int, end int) ([]byte, error) {

	URL := fmt.Sprintf(WPRepoURL, repo, path)
	diffs := fmt.Sprintf("%d:%d", start, end)
	args := []string{URL, "-r", diffs, "--xml", "--summarize"}

	out, err := diff(args...)

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

	out, err := cat(args...)

	return out, err

}

// DoExport writes the remote files from path to local dest.
func DoExport(repo string, path string, dest string) error {

	URL := fmt.Sprintf(WPRepoURL, repo, path)
	args := []string{URL, dest, "-q", "--force", "--depth", "infinity"}

	err := export(args...)

	return err

}

// DoCheckout creates a local svn repo at dest from remote at path.
func DoCheckout(repo string, path string, dest string) error {

	URL := fmt.Sprintf(WPRepoURL, repo, path)
	args := []string{URL, dest, "-q"}

	err := checkout(args...)

	return err

}

// LogResponse contains the response from the `svn log` command.
type LogResponse struct {
	LogEntries []LogEntry `xml:"logentry"`
}

type LogEntry struct {
	Revision int    `xml:"revision,attr"`
	Author   string `xml:"author"`
	Date     string `xml:"date"`
	Paths    []Path `xml:"paths>path"`
	Msg      string `xml:"msg"`
}

type Path struct {
	Kind      string `xml:"kind,attr"`
	TextMods  string `xml:"text-mods,attr"`
	PropsMods string `xml:"prop-mods,attr"`
	Action    string `xml:"action,attr"`
	File      string `xml:",chardata"`
}

// log performs the `svn log` command.
func log(args ...string) ([]LogEntry, error) {
	// Force xml format as it is required below.
	args = append([]string{"log", "--xml"}, args...)

	out, err := command(args...)
	if err != nil {
		return []LogEntry{}, errors.New("SVN Command Failed: " + err.Error())
	}

	var response LogResponse
	err = xml.Unmarshal(out, &response)
	if err != nil {
		return []LogEntry{}, errors.New("Cannot Unmarshal XML: " + err.Error())
	}

	return response.LogEntries, nil
}

// ListResponse contains the response from the `svn list` command.
type ListResponse struct {
	List []ListEntry `xml:"list>entry"`
	Path string      `xml:"path,attr"`
}

type ListEntry struct {
	Kind    string `xml:"kind,attr"`
	Name    string `xml:"name"`
	Commits Commit `xml:"commit"`
}

type Commit struct {
	Revision string `xml:"revision,attr"`
	Author   string `xml:"author"`
	Date     string `xml:"date"`
}

// list performs the `svn list` command.
func list(args ...string) ([]ListEntry, error) {
	// Force xml format as it is required below.
	args = append([]string{"list", "--xml"}, args...)

	out, err := command(args...)
	if err != nil {
		return []ListEntry{}, errors.New("SVN List Command Failed: " + err.Error())
	}

	var response ListResponse
	err = xml.Unmarshal(out, &response)
	if err != nil {
		return []ListEntry{}, errors.New("Cannot Unmarshal XML: " + err.Error())
	}

	return response.List, nil
}

// diff performs the `svn diff` command.
func diff(args ...string) ([]byte, error) {
	cmd := append([]string{"diff"}, args...)

	out, err := exec.Command("svn", cmd...).Output()
	if err != nil {
		return nil, errors.New("SVN Command Failed: " + err.Error())
	}

	return out, nil
}

// export performs the `svn export` command.
func export(args ...string) error {
	args = append([]string{"export"}, args...)

	out, err := command(args...)
	if err != nil {
		argString := strings.Join(args, " ")
		return errors.New("SVN Export Command Failed: " + argString + ": " + err.Error() + ": " + string(out))
	}

	return err
}

// cat performs the `svn cat` command.
func cat(args ...string) ([]byte, error) {
	// Force xml format as it is required below.
	args = append([]string{"cat"}, args...)

	out, err := command(args...)
	if err != nil {
		return nil, errors.New("SVN Cat Command Failed: " + err.Error())
	}

	return out, nil
}

// checkout performs the `svn checkout` command.
func checkout(args ...string) error {
	args = append([]string{"checkout"}, args...)

	_, err := command(args...)
	if err != nil {
		return errors.New("SVN Checkout Command Failed: " + err.Error())
	}

	return err
}

// command executes svn commands returning both output and errors.
func command(args ...string) ([]byte, error) {
	out, err := exec.Command("svn", args...).Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

// IsFolder wraps the `list` command for a specific purpose.
func IsFolder(directory string, path string) bool {
	URL := fmt.Sprintf(WPRepoURL, directory, path)

	args := []string{"list", URL, "--depth=empty"}

	_, err := exec.Command("svn", args...).Output()
	if err != nil {
		return false
	}

	return true
}

// IsClientInstalled checks if the SVN cli client `svn` is available.
func IsClientInstalled() bool {
	_, err := exec.LookPath("svn")
	if err != nil {
		return false
	}

	return true
}
