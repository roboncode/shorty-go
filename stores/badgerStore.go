package stores

import (
	"roboncode.com/go-urlshortener/models"
)

/*

func setupBadger() {
	opts := badger.DefaultOptions
	opts.Dir = "./data/badger"
	opts.ValueDir = "./data/badger"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("answer"), []byte(`{"num":6.13,"strs":["a","b"]}`))
		return err
	})

	err = db.View(func(txn *badger.Txn) error {
		var err error
		var item *badger.Item
		item, err = txn.Get([]byte("answer"))
		if err != nil {
			log.Fatal(err)
		}

		var dat = new(struct {
			Num  float32  `json:"num,omitempty"`
			Strs []string `json:"strs,omitempty"`
		})
		valCopy, _ := item.ValueCopy(nil)
		if err := json.Unmarshal(valCopy, &dat); err != nil {
			panic(err)
		}
		fmt.Println(dat.Num, dat.Strs[1])

		//err = db.View(func(txn *badger.Txn) error {
		//	item, _ := txn.Get([]byte("answer"))
		//	valCopy, _ := item.ValueCopy(nil)
		//	fmt.Printf("The answer is: %s\n", valCopy)
		//	return nil
		//})

		return nil
	})
}
*/

type BadgerStore struct {
}

func (b *BadgerStore) IncCount() int {
	return 0
}

func (b *BadgerStore) Create(code string, url string) (*models.Link, error) {
	return nil, nil
}

func (b *BadgerStore) Read(code string) (*models.Link, error) {
	return nil, nil
}

func (b *BadgerStore) List(limit int64, skip int64) []models.Link {
	return nil
}

func (b *BadgerStore) Delete(code string) int64 {
	return 0
}
