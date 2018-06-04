package repo

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wpdirectory/wpdir/internal/svn"
)

var (
	regexStableTag = regexp.MustCompile(`Stable tag: (.+)`)
)

// PluginReadme ...
type PluginReadme struct {
	Contributors    []string
	Tags            []string
	RequiresAtLeast string
	TestedUpTo      string
	RequiresPHP     string
	StableTag       string
	License         string
	LicenseURI      string
}

// PluginNeedsUpdate ...
func PluginNeedsUpdate(plugin svn.LogEntry) (bool, string) {

	// Get the plugin slug
	parts := strings.Split(plugin.Paths[0].File, "/")
	slug := parts[1]

	// If it is an automated change, ignore it.
	// Mostly commonly the initial empty commit when plugins are approved.
	if plugin.Author == pluginManagementUser {
		return false, ""
	}

	// If all updates were inside the `/assets/` and `/branches/` folders, ignore it.
	// This means that no live code updates occurred.
	var ignoredFolders = true
	for _, commit := range plugin.Paths {

		parts := strings.Split(commit.File, "/")
		if len(parts) > 2 && parts[2] != "assets" && parts[2] != "branches" {
			ignoredFolders = false
		}

	}
	if ignoredFolders == true {
		return false, ""
	}

	stableTag := GetPluginStableTag(plugin)

	var readmeTouched = false
	var tagsTouched []string
	var codeTouched = false

	// TODO: Start by listing all tags which were changed, including trunk as a tag.
	// Loop through changes and look for specific cases
	for _, commit := range plugin.Paths {

		parts := strings.Split(commit.File, "/")

		filename := filepath.Base(commit.File)
		if filename == "readme.txt" || filename == "readme.md" {
			readmeTouched = true
		}

		if len(parts) >= 4 && parts[2] == "tags" && parts[3] != "" {
			tagsTouched = append(tagsTouched, parts[3])
		}

		// Has code been updated
		// Forward slash is included as it might be a folder copy, like tagging a release.
		if commit.File[1:] == "/" || commit.File[4:] == ".php" {
			codeTouched = true
		}

	}

	if readmeTouched == true {

	}

	if len(tagsTouched) > 0 {

		for _, tag := range tagsTouched {

			if tag == stableTag {

				return true, fmt.Sprintf(wpRepoURL, "plugins", slug+"/tags/"+stableTag)

				// TODO: Stable tag is only used if the tags readme and plugin file have the same value.

			}

		}

	}

	if codeTouched == true {

	}

	// Then calculate which tag is currently live, including trunk as a tag.

	// If the live tag was not changed, then we can return false.

	tags, _ := GetPluginTags(slug)
	var tagExists = false
	for _, v := range tags {
		if v == stableTag {
			tagExists = true
		}
	}
	if tagExists != false {
		return true, fmt.Sprintf(wpRepoURL, "plugins", slug+"/tags/"+stableTag)
	}

	// If we have not found a reason to ignore it, assume we should update.
	return true, slug + "/trunk/"

}

// GetPluginStableTag ...
func GetPluginStableTag(plugin svn.LogEntry) string {

	// Get the plugin slug
	parts := strings.Split(plugin.Paths[0].File, "/")
	slug := parts[1]

	out, err := GetCat("plugins", slug, plugin.Revision)
	if err != nil {
		return ""
	}

	// TODO: Currently using fragile regex which captures everything,
	// so 1a.43355674.b0 would be a valid stable tag. Needs improvement,
	// check out the wp.org repo parses them.
	matches := regexStableTag.FindAllStringSubmatch(string(out), -1)

	var stableTag = matches[0][1]

	// If no Stable tag is set fallback to trunk and return.
	if stableTag == "" {
		stableTag = "trunk"
		return stableTag
	}

	// Check if the Stable tag exists in the repo, otherwise fallback to trunk.
	_, err = GetList("plugins", slug+"/tags/"+stableTag)
	if err != nil {
		stableTag = "trunk"
		return stableTag
	}

	return stableTag

}

// GetPluginTags ...
func GetPluginTags(slug string) ([]string, error) {

	out, err := GetList("plugins", slug+"/tags/")
	if err != nil {
		return []string{}, err
	}

	var tags []string

	for _, v := range out {
		tags = append(tags, v.Name)
	}

	return tags, nil

}
