package book

import "context"

// Service holds book business logic.
type Service struct {
	repo Repository
}

// NewService returns a Service backed by repo.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// List returns books with pagination applied.
func (s *Service) List(ctx context.Context, limit, offset int) ([]Book, error) {
	return s.repo.List(ctx, limit, offset)
}

// Get returns a single book or ErrNotFound.
func (s *Service) Get(ctx context.Context, id uint) (Book, error) {
	return s.repo.Get(ctx, id)
}

// Create normalizes, validates, and persists a new book.
func (s *Service) Create(ctx context.Context, b *Book) error {
	b.Normalize()
	if err := b.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, b)
}

// Update applies input to the existing book identified by id.
// Returns ErrNotFound if it does not exist or ValidationError if invalid.
func (s *Service) Update(ctx context.Context, id uint, input Book) (Book, error) {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return Book{}, err
	}

	existing.Title = input.Title
	existing.Author = input.Author
	existing.Rating = input.Rating
	existing.Normalize()
	if err := existing.Validate(); err != nil {
		return Book{}, err
	}

	if err := s.repo.Update(ctx, &existing); err != nil {
		return Book{}, err
	}
	return existing, nil
}

// Delete removes the book identified by id, or returns ErrNotFound.
func (s *Service) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, &existing)
}
