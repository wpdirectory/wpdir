package stats

import (
	"sync"
)

// PluginList ...
type PluginList struct {
	Plugins []Plugin
	sync.RWMutex
}

// Plugin ...
type Plugin struct {
	Slug     string `json:"slug"`
	Revision int    `json:"revision"`
	Date     string `json:"date"`
}

// NewPluginList ...
func NewPluginList() *PluginList {

	plugins := make([]Plugin, 0, 15)
	return &PluginList{Plugins: plugins}

}

// Add adds a Plugin to the list
func (pl *PluginList) Add(plugin Plugin) {

	pl.Lock()
	// Reverse append, prepend.
	pl.Plugins = append([]Plugin{plugin}, pl.Plugins...)
	// Keep list at no larger than 10 items.
	if len(pl.Plugins) > 10 {
		pl.Plugins = pl.Plugins[:len(pl.Plugins)-1]
	}
	pl.Unlock()

}

// GetList returns the plugin list.
func (pl *PluginList) GetList() []Plugin {

	pl.RLock()
	defer pl.RUnlock()
	return pl.Plugins

}

// GetLatestPlugins returns a list of the last updated plugins
func GetLatestPlugins() []Plugin {
	return stats.LatestPlugins.GetList()
}

// AddLatestPlugin ...
func AddLatestPlugin(slug string, revision int, date string) {

	plugin := Plugin{
		Slug:     slug,
		Revision: revision,
		Date:     date,
	}
	stats.LatestPlugins.Add(plugin)

}
