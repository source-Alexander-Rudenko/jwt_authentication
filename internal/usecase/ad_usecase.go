package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"jwt_auth_project/internal/domain"
	"jwt_auth_project/internal/repo"
)

// AdsUseCase описывает бизнес-логику по работе с объявлениями
// CRUD операций и взаимодействие с S3
type AdsUseCase interface {
	InitBucket(ctx context.Context) error
	CreateAd(ctx context.Context, p domain.CreateAdPayload) (*domain.Ad, error)
	GetAdByID(ctx context.Context, id uuid.UUID) (*domain.Ad, error)
	ListAds(ctx context.Context, opts domain.AdListOptions) ([]*domain.Ad, error)
	UpdateAd(ctx context.Context, p domain.UpdateAdPayload) (*domain.Ad, error)
	DeleteAd(ctx context.Context, id uuid.UUID) error
}

// adsUseCase — реализация AdsUseCase
type adsUseCase struct {
	repo         repo.AdsRepository
	s3           *s3.S3
	bucket       string
	validate     *validator.Validate
	maxImageSize int64
}

// NewAdsUsecase создаёт новый экземпляр usecase
// s3Client — клиент из config.NewS3Client(), bucket — название бакета
// maxImageSize — максимальный размер картинки в байтах
func NewAdsUsecase(
	repo repo.AdsRepository,
	s3Client *s3.S3,
	bucket string,
	maxImageSize int64,
) AdsUseCase {
	return &adsUseCase{
		repo:         repo,
		s3:           s3Client,
		bucket:       bucket,
		validate:     validator.New(),
		maxImageSize: maxImageSize,
	}
}

// InitBucket создаёт S3-бакет, если он не существует
func (u *adsUseCase) InitBucket(ctx context.Context) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(u.bucket),
	}
	_, err := u.s3.CreateBucketWithContext(ctx, input)
	return err
}

// CreateAd валидирует payload, загружает картинку в S3 и сохраняет объявление
func (u *adsUseCase) CreateAd(ctx context.Context, p domain.CreateAdPayload) (*domain.Ad, error) {
	p.Title = strings.TrimSpace(p.Title)
	p.Description = strings.TrimSpace(p.Description)
	if err := u.validate.Struct(p); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	if p.ImageSize > u.maxImageSize {
		return nil, errors.New("image file too large")
	}
	id := uuid.New()
	now := time.Now().UTC()

	ext := filepath.Ext(p.ImageName)
	key := fmt.Sprintf("ads/%s%s", id.String(), ext)

	putInput := &s3.PutObjectInput{
		Bucket:        aws.String(u.bucket),
		Key:           aws.String(key),
		Body:          p.Image.(io.ReadSeeker),
		ContentType:   aws.String(p.ContentType),
		ContentLength: aws.Int64(p.ImageSize),
	}
	if _, err := u.s3.PutObjectWithContext(ctx, putInput); err != nil {
		return nil, fmt.Errorf("s3 upload failed: %w", err)
	}
	// 6. Сохранение модели в БД
	ad := &domain.Ad{
		ID:          id,
		AuthorID:    p.AuthorID,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		ImageKey:    key,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := u.repo.CreateAd(ctx, ad); err != nil {
		return nil, fmt.Errorf("db insert failed: %w", err)
	}
	return ad, nil
}

// GetAdByID возвращает объявление по UUID
func (u *adsUseCase) GetAdByID(ctx context.Context, id uuid.UUID) (*domain.Ad, error) {
	return u.repo.GetAdByID(ctx, id)
}

// ListAds возвращает список объявлений
func (u *adsUseCase) ListAds(ctx context.Context, opts domain.AdListOptions) ([]*domain.Ad, error) {
	return u.repo.ListAds(ctx, opts)
}

// UpdateAd обновляет объявление и при необходимости заменяет картинку
func (u *adsUseCase) UpdateAd(ctx context.Context, p domain.UpdateAdPayload) (*domain.Ad, error) {
	if err := u.validate.Struct(p); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	existing, err := u.repo.GetAdByID(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	if p.Image != nil {
		if p.ImageSize > u.maxImageSize {
			return nil, errors.New("image file too large")
		}
		ext := filepath.Ext(p.ImageName)
		key := fmt.Sprintf("ads/%s%s", existing.ID.String(), ext)
		putInput := &s3.PutObjectInput{
			Bucket:        aws.String(u.bucket),
			Key:           aws.String(key),
			Body:          p.Image.(io.ReadSeeker),
			ContentType:   aws.String(p.ContentType),
			ContentLength: aws.Int64(p.ImageSize),
		}
		if _, err := u.s3.PutObjectWithContext(ctx, putInput); err != nil {
			return nil, fmt.Errorf("s3 upload failed: %w", err)
		}
		existing.ImageKey = key
	}

	existing.Title = p.Title
	existing.Description = p.Description
	existing.Price = p.Price
	existing.UpdatedAt = time.Now().UTC()

	if err := u.repo.UpdateAd(ctx, existing); err != nil {
		return nil, fmt.Errorf("db update failed: %w", err)
	}
	return existing, nil
}

func (u *adsUseCase) DeleteAd(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteAd(ctx, id)
}
