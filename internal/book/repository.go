package book

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository abstracts book persistence.
type Repository interface {
	List(ctx context.Context, limit, offset int) ([]Book, error)
	Get(ctx context.Context, id uint) (Book, error)
	Create(ctx context.Context, b *Book) error
	Update(ctx context.Context, b *Book) error
	Delete(ctx context.Context, b *Book) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository returns a GORM-backed Repository.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) List(ctx context.Context, limit, offset int) ([]Book, error) {
	var books []Book
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&books).Error
	return books, err
}

func (r *gormRepository) Get(ctx context.Context, id uint) (Book, error) {
	var b Book
	err := r.db.WithContext(ctx).First(&b, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Book{}, ErrNotFound
	}
	return b, err
}

func (r *gormRepository) Create(ctx context.Context, b *Book) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *gormRepository) Update(ctx context.Context, b *Book) error {
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *gormRepository) Delete(ctx context.Context, b *Book) error {
	return r.db.WithContext(ctx).Delete(b).Error
}
