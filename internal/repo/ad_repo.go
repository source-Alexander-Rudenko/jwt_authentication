package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"jwt_auth_project/internal/domain"
)

var ErrAdNotFound = errors.New("ad not found")

type AdsRepo struct {
	pool *pgxpool.Pool
}

func NewAdsRepo(pool *pgxpool.Pool) *AdsRepo {
	return &AdsRepo{pool: pool}
}

type AdsRepository interface {
	CreateAd(ctx context.Context, ad *domain.Ad) error
	GetAdByID(ctx context.Context, id uuid.UUID) (*domain.Ad, error)
	ListAds(ctx context.Context, opts domain.AdListOptions) ([]*domain.Ad, error)
	UpdateAd(ctx context.Context, ad *domain.Ad) error
	DeleteAd(ctx context.Context, id uuid.UUID) error
}

func (r *AdsRepo) CreateAd(ctx context.Context, ad *domain.Ad) error {
	_, err := r.pool.Exec(ctx, `
        INSERT INTO "ADS" (
            id, author_id, title, description, price, image_key, created_at, updated_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
    `, ad.ID, ad.AuthorID, ad.Title, ad.Description, ad.Price, ad.ImageKey, ad.CreatedAt, ad.UpdatedAt)
	return err
}

func (r *AdsRepo) GetAdByID(ctx context.Context, id uuid.UUID) (*domain.Ad, error) {
	a := new(domain.Ad)
	err := r.pool.QueryRow(ctx, `
        SELECT id, author_id, title, description, price, image_key, created_at, updated_at
        FROM "ADS"
        WHERE id = $1
    `, id).Scan(
		&a.ID, &a.AuthorID, &a.Title, &a.Description,
		&a.Price, &a.ImageKey, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAdNotFound
		}
		return nil, err
	}
	return a, nil
}

func (r *AdsRepo) ListAds(ctx context.Context, opts domain.AdListOptions) ([]*domain.Ad, error) {
	var field string
	switch opts.SortField {
	case "price":
		field = "price"
	case "created_at":
		field = "created_at"
	default:
		field = "created_at"
	}

	dir := "DESC"
	if opts.SortAsc {
		dir = "ASC"
	}

	sb := squirrel.
		Select(
			"id",
			"author_id",
			"title",
			"description",
			"price",
			"image_key",
			"created_at",
			"updated_at",
		).
		From(`"ADS"`).
		OrderBy(fmt.Sprintf("%s %s", field, dir)).
		Limit(uint64(opts.Limit)).
		Offset(uint64(opts.Offset)).
		PlaceholderFormat(squirrel.Dollar)

	sqlStr, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Ad
	for rows.Next() {
		a := new(domain.Ad)
		if err := rows.Scan(
			&a.ID,
			&a.AuthorID,
			&a.Title,
			&a.Description,
			&a.Price,
			&a.ImageKey,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *AdsRepo) UpdateAd(ctx context.Context, ad *domain.Ad) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE "ADS"
        SET title = $2, description = $3, price = $4, image_key = $5, updated_at = $6
        WHERE id = $1
    `, ad.ID, ad.Title, ad.Description, ad.Price, ad.ImageKey, ad.UpdatedAt)
	return err
}

func (r *AdsRepo) DeleteAd(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM "ADS" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrAdNotFound
	}
	return nil
}
