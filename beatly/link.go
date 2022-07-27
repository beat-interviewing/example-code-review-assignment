package beatly

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
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

	// Visits holds past visits for analytics purposes. See the VisitsPer method
	// on how visits can be aggregated per second, minute, hour or day.
	Visits []time.Time
}

// VisitsPer aggregates the number of visits by the selected interval. The
// result is a map where keys are times and values are the number of visits.
func (link *Link) VisitsPer(interval time.Duration) (visits map[string]int) {

	var key string

	switch interval {
	case time.Second:
		key = "2006-01-02T15:04:05"
	case time.Minute:
		key = "2006-01-02T15:04"
	case time.Hour:
		key = "2006-01-02T15"
	case time.Hour * 24:
		key = "2006-01-02"
	default:
		panic("invalid argument")
	}

	visits = make(map[string]int)

	// Count visits grouped by the chosen interval.
	for _, visit := range link.Visits {

		visitKey := visit.Format(key)

		_, ok := visits[visitKey]
		if !ok {
			visits[visitKey] = 0
		}
		visits[visitKey]++
	}

	return visits
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
