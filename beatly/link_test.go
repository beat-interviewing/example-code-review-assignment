package beatly

import (
	"testing"
	"time"
)

func TestVisitsPer(t *testing.T) {

	now := time.Now()

	for _, test := range []struct {
		link Link
		len  int
	}{
		{
			Link{Visits: []time.Time{}},
			0,
		},
		{
			Link{Visits: []time.Time{now}},
			1,
		},
		{
			Link{
				Visits: []time.Time{
					now,
					now.Add(1 * time.Second),
					now.Add(1 * time.Second),
					now.Add(2 * time.Second),
					now.Add(2 * time.Second),
					now.Add(2 * time.Second),
				},
			},
			3,
		},
	} {
		t.Run("1s", func(t *testing.T) {
			visits := test.link.VisitsPer(time.Second)
			if len(visits) != test.len {
				t.Errorf("unexpected len(visits): expected %v, have %v", test.len, len(visits))
			}
		})
	}
}
