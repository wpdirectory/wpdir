package db

import (
	"log"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db      *bolt.DB
	buckets = []string{
		"searches",
		"plugins",
		"themes",
	}
)

// Close closes the bolt db.
// Should be called with defer after init in the main() func.
func Close() {
	err := db.Close()
	if err != nil {
		log.Println(err)
	}
}

// New opens a new bolt db and ensured default buckets exist.
func New(dataDir string) {

	path := filepath.Join(dataDir, "wpdir.db")
	options := &bolt.Options{
		Timeout: 1 * time.Second,
	}

	var err error
	db, err = bolt.Open(path, 0600, options)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure main Buckets exist
	err = db.Update(func(tx *bolt.Tx) error {

		for _, name := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalln("Cannot create default buckets: ", err)
	}

}

// SaveItemToBucket ...
func SaveItemToBucket(key string, content []byte, bucket string) error {

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("searches"))
		err := b.Put([]byte(key), content)
		return err
	})

	return err

}
