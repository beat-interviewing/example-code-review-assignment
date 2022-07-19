package beatly

import (
	"fmt"
	"net/http"
	"time"
)

// Service is a core building block of the BEAT.ly service. It provides methods
// which can create, read and visit shortened links.
type Service interface {

	// Create creates a link.
	//
	// The provided link will be validated and an error is returned if
	// validation fails.
	//
	// If the operation is successful the provided link will have its ID and
	// IDHash fields set to an auto-incremented id and its hash respectively.
	Create(link *Link) error

	// Read returns a link matching the supplied id. If no link was found error
	// will be non-nil.
	Read(id string) (*Link, error)

	// Visit retrieves a link matching the provided id, and registers the time
	// of retrieval for analytics purposes.
	Visit(id string) (*Link, error)
}

type service struct {
	store Store
}

// NewService returns an implementation of Service.
func NewService(s Store) Service {
	return &service{s}
}

func (s *service) Create(link *Link) error {

	// The redirect field is optional. If left unset we'll default to 302 Found.
	if link.Redirect == 0 {
		link.Redirect = http.StatusFound
	}

	if err := link.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return s.store.Create(link)
}

func (s *service) Read(id string) (*Link, error) {
	return s.store.Read(id)
}

func (s *service) Visit(id string) (*Link, error) {
	return s.store.Visit(id, time.Now())
}