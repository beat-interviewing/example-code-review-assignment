package beatly

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/speps/go-hashids"
	bolt "go.etcd.io/bbolt"
)

// Store is the persistence layer of the BEAT.ly service.
type Store interface {

	// Create persists a link to disk.
	Create(*Link) error

	// Read retrieves a link from disk matching the provided id.
	Read(id string) (*Link, error)

	// Visit retrieves a link from disk matching the provided id and registers
	// the time of each request for analytics purposes.
	Visit(id string) (*Link, error)
}

type store struct {
	db *bolt.DB
}

// NewBoltStore returns an implementation of Store which relies on
// BoltDB for persistence.
//
// Bolt is a pure Go key/value store with a goal to provide a simple, fast, and
// reliable database for projects that don't require a full database server.
//
// As such programs that rely on this implementation of Store cannot be
// deployed in a replicated manner.
func NewBoltStore(dsn string) (Store, error) {
	db, err := bolt.Open(dsn, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &store{db}, nil
}

func (s *store) Create(link *Link) error {

	return s.db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte("links"))
		if err != nil {
			return err
		}

		id, _ := b.NextSequence()

		link.ID = int(id)
		link.IDHash = hash(link.ID)

		buf, err := json.Marshal(link)
		if err != nil {
			return err
		}

		return b.Put([]byte(link.IDHash), buf)
	})
}

func (s *store) Read(id string) (link *Link, err error) {

	err = s.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("links"))
		if b == nil {
			return fmt.Errorf("not found")
		}

		buf := b.Get([]byte(id))
		if buf == nil {
			return fmt.Errorf("not found")
		}

		return json.Unmarshal(buf, &link)
	})

	return
}

func (s *store) Visit(id string) (link *Link, err error) {

	err = s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("links"))
		if b == nil {
			return fmt.Errorf("not found")
		}

		buf := b.Get([]byte(id))
		if buf == nil {
			return fmt.Errorf("not found")
		}

		err := json.Unmarshal(buf, &link)
		if err != nil {
			return err
		}

		link.Visits = append(link.Visits, time.Now())

		buf, err = json.Marshal(link)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), buf)
	})

	return
}

var hd = &hashids.HashIDData{
	Salt:      "beat.ly",
	MinLength: 3,
	Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890",
}

var h, _ = hashids.NewWithData(hd)

func hash(in int) string {
	s, err := h.Encode([]int{in})
	if err != nil {
		panic(err)
	}
	return s
}
