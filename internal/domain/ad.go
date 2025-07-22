package domain

import (
	"github.com/google/uuid"
	"io"
	"time"
)

type Ad struct {
	ID          uuid.UUID `json:"id"`
	AuthorID    uuid.UUID `json:"author_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageKey    string    `json:"image_key"` // ключ в S3
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdListOptions struct {
	Limit     int
	Offset    int
	SortField string // "price" или "created_at"
	SortAsc   bool   // true = ASC, false = DESC
	MinPrice  *float64
	MaxPrice  *float64
}

type CreateAdPayload struct {
	AuthorID    uuid.UUID `json:"author_id" validate:"required"`
	Title       string    `json:"title"       validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"required,min=10,max=1000"`
	Price       float64   `json:"price"       validate:"required,gte=0"`
	Image       io.Reader `json:"-"           validate:"required"`
	ImageSize   int64     `json:"-"           validate:"required,gte=1"`
	ImageName   string    `json:"-"           validate:"required"`
	ContentType string    `json:"-"           validate:"required"`
}

type UpdateAdPayload struct {
	ID          uuid.UUID `json:"id"          validate:"required"`
	Title       string    `json:"title"       validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"required,min=10,max=1000"`
	Price       float64   `json:"price"       validate:"required,gte=0"`
	Image       io.Reader `json:"-"           validate:"omitempty"`
	ImageSize   int64     `json:"-"           validate:"omitempty,gte=1"`
	ImageName   string    `json:"-"           validate:"omitempty"`
	ContentType string    `json:"-"           validate:"omitempty"`
}
