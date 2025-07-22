package config

import (
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Config хранит параметры для подключения к S3
type S3Config struct {
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
}

// NewConfigFromEnv читает все параметры из переменных окружения
func NewConfigFromEnv() *S3Config {
	return &S3Config{
		Bucket:    os.Getenv("S3_BUCKET"),
		Region:    os.Getenv("S3_REGION"),
		Endpoint:  os.Getenv("S3_ENDPOINT"),
		AccessKey: os.Getenv("S3_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
	}
}

// NewS3Client создаёт и возвращает *s3.S3
func NewS3Client(cfg *S3Config) *s3.S3 {
	awsCfg := &aws.Config{
		Region: aws.String(cfg.Region),
	}

	// если задан endpoint (например, MinIO), подменяем
	if cfg.Endpoint != "" {
		awsCfg.Endpoint = aws.String(cfg.Endpoint)
		awsCfg.S3ForcePathStyle = aws.Bool(true)
	}

	// используем статические креды из env
	awsCfg.Credentials = credentials.NewStaticCredentials(
		cfg.AccessKey,
		cfg.SecretKey,
		"",
	)

	sess, err := session.NewSession(awsCfg)
	if err != nil {
		log.Fatalf("failed to create AWS session: %v", err)
	}
	slog.Info("s3 client created")
	return s3.New(sess)
}
