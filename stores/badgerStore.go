package stores

import (
	"github.com/dgraph-io/badger"
	"github.com/roboncode/go-urlshortener/helpers"
	"github.com/roboncode/go-urlshortener/models"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	BadgerDir = "BADGER_DIR"
)

type BadgerStore struct {
	codeDB    *badger.DB
	hashDB    *badger.DB
	counterDB *badger.DB
	counter   models.Counter
}

func NewBadgerStore() Store {
	viper.SetDefault(BadgerDir, "./data/badger")
	_ = viper.BindEnv(BadgerDir)

	b := BadgerStore{}
	b.codeDB = b.connect("code")
	b.hashDB = b.connect("hash")
	b.counterDB = b.connect("counter")
	b.restoreCounter()
	return &b
}

func (b *BadgerStore) connect(database string) *badger.DB {
	baseDir := viper.GetString(BadgerDir)
	dir := baseDir + "/" + database
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
	_ = b.counterDB.View(func(txn *badger.Txn) error {
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

	err = b.codeDB.View(func(txn *badger.Txn) error {
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

func (b *BadgerStore) IncCount() int {
	b.counter.Value++
	_ = b.counterDB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("counter"), b.counter.EncodeCounter())
	})
	return b.counter.Value
}

func (b *BadgerStore) Create(code string, longUrl string) (*models.Link, error) {
	var err error
	hashedLongUrl := helpers.MD5(longUrl)
	link := b.FindLink(hashedLongUrl)
	if link != nil {
		link.ShortUrl = helpers.GetShortUrl(link.Code)
		return link, err
	}

	newLink := models.Link{
		LongUrl: longUrl,
		Code:    code,
		Created: time.Now(),
	}

	encodedLink := newLink.EncodeLink()

	err = b.codeDB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(code), encodedLink)
	})

	err = b.hashDB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(hashedLongUrl), encodedLink)
	})

	if err != nil {
		return nil, err
	}

	newLink.ShortUrl = helpers.GetShortUrl(newLink.Code)
	return &newLink, nil
}

func (b *BadgerStore) Read(code string) (*models.Link, error) {
	var link *models.Link
	var err error

	err = b.codeDB.View(func(txn *badger.Txn) error {
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
	link.ShortUrl = helpers.GetShortUrl(link.Code)
	return link, nil
}

func (b *BadgerStore) List(limit int, skip int) []models.Link {
	var link *models.Link
	links := make([]models.Link, 0)

	err := b.codeDB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = limit
		it := txn.NewIterator(opts)
		defer it.Close()
		//prefix := []byte("k1pGP")
		prefix := []byte("")
		//for it.Rewind(); it.Valid(); it.Next() {
		//for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		for it.Seek(prefix); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, _ := item.Value()

			//fmt.Printf("key=%s, value=%s\n", k, v)

			link, _ = models.DecodeLink(v)
			link.ID = k
			link.ShortUrl = helpers.GetShortUrl(link.Code)
			links = append(links, *link)
		}
		return nil
	})

	if err != nil {
		return links
	}

	return links
}

func (b *BadgerStore) Delete(code string) int {
	link, _ := b.Read(code)
	var err error
	if link != nil {
		err = b.codeDB.Update(func(txn *badger.Txn) error {
			err = txn.Delete([]byte(code))
			return err
		})
		if err != nil {
			return 0
		}
		err = b.hashDB.Update(func(txn *badger.Txn) error {
			hashedLongUrl := helpers.MD5(link.LongUrl)
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
