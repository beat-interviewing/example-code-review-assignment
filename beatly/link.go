package beatly

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Link struct {
	ID       int
	IDHash   string
	Target   string
	Redirect int
	Visits   []time.Time
}

// VisitsPer counts visits grouped by granularity.
func (link *Link) VisitsPer(granularity time.Duration) (visits map[string]int) {

	var key string

	switch granularity {
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

	// Count visits grouped by the chosen granularity.
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
