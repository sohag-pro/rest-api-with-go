package book

import (
	"context"
	"errors"
	"testing"
)

// fakeRepo is an in-memory Repository for unit-testing the service.
type fakeRepo struct {
	books  map[uint]Book
	nextID uint
}

func newFakeRepo() *fakeRepo { return &fakeRepo{books: map[uint]Book{}, nextID: 1} }

func (f *fakeRepo) List(_ context.Context, limit, offset int) ([]Book, error) {
	var out []Book
	for _, b := range f.books {
		out = append(out, b)
	}
	if offset > len(out) {
		return nil, nil
	}
	out = out[offset:]
	if limit < len(out) {
		out = out[:limit]
	}
	return out, nil
}

func (f *fakeRepo) Get(_ context.Context, id uint) (Book, error) {
	b, ok := f.books[id]
	if !ok {
		return Book{}, ErrNotFound
	}
	return b, nil
}

func (f *fakeRepo) Create(_ context.Context, b *Book) error {
	b.ID = f.nextID
	f.nextID++
	f.books[b.ID] = *b
	return nil
}

func (f *fakeRepo) Update(_ context.Context, b *Book) error {
	f.books[b.ID] = *b
	return nil
}

func (f *fakeRepo) Delete(_ context.Context, b *Book) error {
	delete(f.books, b.ID)
	return nil
}

func TestServiceCreateValidates(t *testing.T) {
	svc := NewService(newFakeRepo())
	ctx := context.Background()

	if err := svc.Create(ctx, &Book{Title: "  ", Rating: 1}); err == nil {
		t.Fatal("expected validation error for empty title")
	}
	if err := svc.Create(ctx, &Book{Title: "Go", Rating: 9}); err == nil {
		t.Fatal("expected validation error for bad rating")
	}

	b := &Book{Title: "  Go  ", Author: " Pike ", Rating: 5}
	if err := svc.Create(ctx, b); err != nil {
		t.Fatalf("create: %v", err)
	}
	if b.Title != "Go" || b.Author != "Pike" {
		t.Fatalf("normalize failed: %+v", b)
	}
}

func TestServiceUpdateNotFound(t *testing.T) {
	svc := NewService(newFakeRepo())
	_, err := svc.Update(context.Background(), 42, Book{Title: "X", Rating: 1})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestServiceDeleteNotFound(t *testing.T) {
	svc := NewService(newFakeRepo())
	if err := svc.Delete(context.Background(), 42); !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
