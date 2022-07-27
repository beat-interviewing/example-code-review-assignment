package beatly

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/speps/go-hashids"
)

// Store is the persistence layer of the BEAT.ly service.
type Store interface {

	// Create persists a link to disk. Upon successful completion the link will
	// have its ID and IDHash fields set.
	Create(link *Link) error

	// Visit retrieves a link from disk matching the provided id hash and
	// registers the time of each request for analytics purposes.
	Visit(hash string) (*Link, error)
}

type sqlite struct {
	db *sql.DB
	m  *migrate.Migrate
}

// NewSQLiteStore returns an implementation of Store which relies on SQLite for
// persistence.
//
// SQLite is a library that implements a small, fast, self-contained,
// high-reliability, full-featured, SQL database engine.
func NewSQLiteStore(dsn string) (Store, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	m, err := migrate.New("file://migrations", fmt.Sprintf("sqlite3://%s", dsn))
	if err != nil {
		return nil, err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	return &sqlite{db, m}, nil
}

func (s *sqlite) Create(link *Link) error {

	r, err := s.db.Exec(`insert into links(target, redirect) values (?, ?)`, link.Target, link.Redirect)
	if err != nil {
		return fmt.Errorf("insert failed. %w", err)
	}

	id, _ := r.LastInsertId()

	link.ID = id
	link.IDHash, err = Encode(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *sqlite) Visit(hash string) (*Link, error) {

	// Decode the hash to get the numeric id of the link. From this moment on we
	// will be mostly using the id.
	id, err := Decode(hash)
	if err != nil {
		return nil, err
	}

	link := &Link{
		ID:     id,
		IDHash: hash,
	}

	// Query the database for the link by its id.
	row := s.db.QueryRow(`select id, target, redirect from links where id = ?`, id)
	err = row.Scan(&link.ID, &link.Target, &link.Redirect)
	if err != nil {
		return nil, err
	}

	return link, nil
}

// Using the hashids library we will generate short, unique, non-sequential
// hashes from numeric ids.
//
// Hashes are preferable to numeric auto-incrementing ids by being less
// predictable or guessable.
//
// The configuration below uses a random 32 character string as salt, making the
// generated hashes harder to guess.
//
// The minimum length of hashes will be at least 3 characters long.
//
// The generated hashes will be a encoded in base62 which is more friendly to
// URLs than the popular base64 encoding as it omits the `+` and `/` characters.
//
// References:
// 	- https://en.wikipedia.org/wiki/URL_shortening#Techniques
//  - https://en.wikipedia.org/wiki/Base62
// 	- https://hashids.org/go/
//
var hd = &hashids.HashIDData{
	Salt:      "f29490b05e6049908ae6aa6d6312ea85",
	MinLength: 3,
	Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890",
}

var h, _ = hashids.NewWithData(hd)

// Encode a numeric id to a hash.
func Encode(i int64) (string, error) {
	return h.EncodeInt64([]int64{i})
}

// Decode reverses a hash back to its original numeric value.
func Decode(s string) (int64, error) {
	i, err := h.DecodeInt64WithError(s)
	if err != nil {
		return 0, err
	}
	return i[0], nil
}
