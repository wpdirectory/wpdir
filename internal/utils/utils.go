package utils

import (
	"io"
	"io/ioutil"
	"net/url"
)

// EncodeURL properly encodes the URL for compatibility with special characters
// e.g. 新浪微博 and ЯндексФотки
func EncodeURL(rawURL string) string {

	u, _ := url.Parse(rawURL)

	URL := u.String()

	return URL

}

// CheckClose is used to check the return from Close in a defer statement.
func CheckClose(c io.Closer, err *error) {
	cErr := c.Close()
	if *err == nil {
		*err = cErr
	}
}

// DrainAndClose discards all data from rd and closes it.
func DrainAndClose(rd io.ReadCloser, err *error) {
	if rd == nil {
		return
	}
	_, _ = io.Copy(ioutil.Discard, rd)
	cErr := rd.Close()
	if err != nil && *err == nil {
		*err = cErr
	}
}
