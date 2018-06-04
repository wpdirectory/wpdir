package stats

import (
	"sync"
)

// ThemeList ...
type ThemeList struct {
	Themes []Theme
	sync.RWMutex
}

// Theme ...
type Theme struct {
	Slug     string `json:"slug"`
	Revision int    `json:"revision"`
	Date     string `json:"date"`
}

// NewThemeList ...
func NewThemeList() *ThemeList {

	themes := make([]Theme, 0, 15)
	return &ThemeList{Themes: themes}

}

// Add adds a Theme to the list
func (th *ThemeList) Add(theme Theme) {

	th.Lock()
	// Reverse append, prepend.
	th.Themes = append([]Theme{theme}, th.Themes...)
	// Keep list at no larger than 10 items.
	if len(th.Themes) > 10 {
		th.Themes = th.Themes[:len(th.Themes)-1]
	}
	th.Unlock()

}

// GetList returns the Theme list.
func (th *ThemeList) GetList() []Theme {

	th.RLock()
	defer th.RUnlock()
	return th.Themes

}

// GetLatestThemes returns a list of the last updated themes
func GetLatestThemes() []Theme {
	return stats.LatestThemes.GetList()
}

// AddLatestTheme ...
func AddLatestTheme(slug string, revision int, date string) {

	theme := Theme{
		Slug:     slug,
		Revision: revision,
		Date:     date,
	}
	stats.LatestThemes.Add(theme)

}
