package stores

import (
	"github.com/dgraph-io/badger"
	"github.com/roboncode/go-urlshortener/helpers"
	"github.com/roboncode/go-urlshortener/models"
	"github.com/spf13/viper"
	"os"
	"time"
)

const BadgerDir = "BADGER_DIR"

type BadgerStore struct {
	db      *badger.DB
	counter models.Counter
}

func NewBadgerStore() *BadgerStore {
	viper.SetDefault(BadgerDir, "./data/badger")
	_ = viper.BindEnv(BadgerDir)

	b := BadgerStore{}
	b.db = b.connect()
	b.restoreCounter()
	return &b
}

func (b *BadgerStore) connect() *badger.DB {
	dir := viper.GetString(BadgerDir)
	_ = os.MkdirAll(dir, os.ModePerm)

	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	return db
}

func (b *BadgerStore) restoreCounter() {
	_ = b.db.View(func(txn *badger.Txn) error {
		var err error
		var item *badger.Item
		item, err = txn.Get([]byte("counter"))
		if err != nil {
			return err
		}

		val, err := item.Value()
		if err != nil {
			return err
		}

		b.counter, err = models.DecodeCounter(val)
		return err
	})
}

func (b *BadgerStore) FindLink(hashedLongUrl string) *models.Link {
	var link *models.Link
	var err error

	err = b.db.View(func(txn *badger.Txn) error {
		var err error
		var item *badger.Item
		item, err = txn.Get([]byte(hashedLongUrl))
		if err != nil {
			return err
		}

		val, err := item.Value()
		if err != nil {
			return err
		}

		link, err = models.DecodeLink(val)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil
	}
	return link
}

func (b *BadgerStore) IncCount() int64 {
	b.counter.Value++
	_ = b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("counter"), b.counter.EncodeCounter())
	})
	return b.counter.Value
}

func (b *BadgerStore) Create(code string, longUrl string) (*models.Link, error) {
	var err error
	hashedLongUrl := helpers.MD5(longUrl)
	link := b.FindLink(hashedLongUrl)
	if link != nil {
		helpers.FormatShortUrl(link)
		return link, err
	}

	newLink := models.Link{
		LongUrl: longUrl,
		Code:    code,
		Created: time.Now(),
	}

	err = b.db.Update(func(txn *badger.Txn) error {
		encodedLink := newLink.EncodeLink()
		err = txn.Set([]byte(code), encodedLink)

		err = txn.Set([]byte(hashedLongUrl), encodedLink)
		return err
	})

	if err != nil {
		return nil, err
	}

	helpers.FormatShortUrl(&newLink)
	return &newLink, nil
}

func (b *BadgerStore) Read(code string) (*models.Link, error) {
	var link *models.Link
	var err error

	err = b.db.View(func(txn *badger.Txn) error {
		var err error
		var item *badger.Item
		item, err = txn.Get([]byte(code))
		if err != nil {
			return err
		}

		val, err := item.Value()
		if err != nil {
			return err
		}

		link, err = models.DecodeLink(val)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	helpers.FormatShortUrl(link)
	return link, nil
}

func (b *BadgerStore) List(limit int64, skip int64) []models.Link {
	return make([]models.Link, 0)
}

func (b *BadgerStore) Delete(code string) int64 {
	link, _ := b.Read(code)
	var err error
	if link != nil {
		err = b.db.Update(func(txn *badger.Txn) error {
			hashedLongUrl := helpers.MD5(link.LongUrl)
			err = txn.Delete([]byte(code))
			err = txn.Delete([]byte(hashedLongUrl))
			return err
		})
		if err != nil {
			return 0
		}
		return 1
	}
	return 0
}
