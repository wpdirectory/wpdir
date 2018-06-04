package stats

import (
	"time"
)

var (
	stats *Stats
)

// Stats holds summary data
type Stats struct {
	LatestPlugins *PluginList
	LatestThemes  *ThemeList
	TotalPlugins  int
	TotalThemes   int
}

// Extension ...
type Extension struct {
	Slug      string
	ActiveTag string
	Revision  int
	Time      time.Time
}

// Setup ...
func Setup() {

	pluginList := NewPluginList()
	themeList := NewThemeList()

	stats = &Stats{
		LatestPlugins: pluginList,
		LatestThemes:  themeList,
	}

}
