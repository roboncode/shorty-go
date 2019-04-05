package stores

import "github.com/roboncode/go-urlshortener/models"

type Store interface {
	IncCount() int64
	Create(code string, url string) (*models.Link, error)
	Read(code string) (*models.Link, error)
	List(limit int64, skip int64) []models.Link
	Delete(code string) int64
}
