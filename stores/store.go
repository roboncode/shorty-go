package stores

import "github.com/roboncode/go-urlshortener/models"

type Store interface {
	IncCount() int
	Create(code string, url string) (*models.Link, error)
	Read(code string) (*models.Link, error)
	List(limit int, skip int) []models.Link
	Delete(code string) int
}
