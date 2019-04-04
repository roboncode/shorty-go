package stores

import (
	"roboncode.com/go-urlshortener/models"
)

type BadgerStore struct {
}

func (b *BadgerStore) IncCount() int {

}

func (b *BadgerStore) Create(code string, url string) (*models.Link, error) {

}

func (b *BadgerStore) Read(code string) (*models.Link, error) {

}

func (b *BadgerStore) List(limit int64, skip int64) []models.Link {

}

func (b *BadgerStore) Delete(code string) int64 {

}
