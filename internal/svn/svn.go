package svn

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	wpRepoURL = "https://%s.svn.wordpress.org/%s"
)

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

// Log performs the `svn log` command.
func Log(args ...string) ([]LogEntry, error) {

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

// List performs the `svn list` command.
func List(args ...string) ([]ListEntry, error) {

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

// Diff performs the `svn diff` command.
func Diff(args ...string) ([]byte, error) {

	cmd := append([]string{"diff"}, args...)

	out, err := exec.Command("svn", cmd...).Output()
	if err != nil {
		return nil, errors.New("SVN Command Failed: " + err.Error())
	}

	return out, nil

}

// Export performs the `svn export` command.
func Export(args ...string) error {

	args = append([]string{"export"}, args...)

	out, err := command(args...)
	if err != nil {
		argString := strings.Join(args, " ")
		return errors.New("SVN command failed: svn export " + argString + ": " + err.Error() + ": " + string(out))
	}

	return err

}

// Cat performs the `svn cat` command.
func Cat(args ...string) ([]byte, error) {

	// Force xml format as it is required below.
	args = append([]string{"cat"}, args...)

	out, err := command(args...)
	if err != nil {
		return nil, errors.New("SVN Cat Command Failed: " + err.Error())
	}

	return out, nil

}

// Checkout performs the `svn checkout` command.
func Checkout(args ...string) error {

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

// IsFolder wraps the `List` func for a specific purpose.
func IsFolder(directory string, path string) bool {

	URL := fmt.Sprintf(wpRepoURL, directory, path)

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
