package beatly

import (
	"fmt"
	"net/http"
	"net/url"
)

// Link is the internal representation of a shortened link.
type Link struct {

	// The link id as determined by the underlying data store. This will likely
	// be an auto-incrementing integer in most relational databases.
	ID int64

	// The link id hash is computed using the link id as input.
	IDHash string

	// Target is a non-shortened URL which the service will redirect to when the
	// short URL is visited.
	Target string

	// Redirect determines the 3xx HTTP status code used when performing the
	// redirect.
	Redirect int
}

func (link *Link) Validate() error {

	_, err := url.ParseRequestURI(link.Target)
	if err != nil {
		return fmt.Errorf("field `target` is invalid: expected a URL. %w", err)
	}

	if link.Redirect != http.StatusMovedPermanently &&
		link.Redirect != http.StatusFound &&
		link.Redirect != http.StatusTemporaryRedirect &&
		link.Redirect != http.StatusPermanentRedirect {
		return fmt.Errorf("field `redirect` is invalid: expected one of (301, 302, 307, 308)")
	}

	return nil
}
