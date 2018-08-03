package db

import (
	"errors"
	"log"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db      *bolt.DB
	buckets = []string{
		"repos",
		"plugins",
		"themes",
		"searches",
	}
	searchBuckets = []string{
		"search_data",
		"all_dates",
		"public_dates",
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

// Setup opens a new bolt db and ensures default buckets exist.
func Setup(dir string) {
	path := filepath.Join(dir, "data", "db", "wpdir.db")
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
			b, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return err
			}
			if name == "searches" {
				for _, name := range searchBuckets {
					_, err := b.CreateBucketIfNotExists([]byte(name))
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln("Cannot create default buckets: ", err)
	}
}

// PutToBucket adds an item to bucket
func PutToBucket(key string, content []byte, bucket string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(key), content)
	})
	return err
}

// GetFromBucket returns an item from a bucket
func GetFromBucket(key string, bucket string) ([]byte, error) {
	var data []byte
	var err error
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data = b.Get([]byte(key))
		return nil
	})
	if len(data) == 0 {
		err = errors.New("No data found")
	}
	return data, err
}

// DeleteFromBucket deletes an item from a bucket
func DeleteFromBucket(key string, bucket string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(key))
	})
	return err
}

// GetAllFromBucket returns all items from a bucket
func GetAllFromBucket(bucket string) (map[string][]byte, error) {
	items := make(map[string][]byte)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.ForEach(func(k, v []byte) error {
			items[string(k)] = v
			return nil
		})
		return err
	})

	return items, err
}

// GetLatestPublicSearchList most recent public searches
func GetLatestPublicSearchList(limit int) []string {
	var list []string
	db.View(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("searches"))
		// Get Relevant Internal Buckets
		dates := s.Bucket([]byte("public_dates")).Cursor()
		if dates == nil {
			return nil
		}
		i := 0
		for k, v := dates.Last(); k != nil; k, v = dates.Prev() {
			list = append(list, string(v))
			i++
			if i == limit {
				break
			}
		}
		return nil
	})
	return list
}

// DeleteSearches removes all Search data and buckets
func DeleteSearches() error {
	// Start the transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.DeleteBucket([]byte("searches")); err != nil {
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// SaveSearch saves the Search data to DB
func SaveSearch(searchID string, created string, private bool, bytes []byte) error {
	// Start the transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get root Searches bucket
	s := tx.Bucket([]byte("searches"))

	// Get Relevant Internal Buckets
	data := s.Bucket([]byte("search_data"))
	allDates := s.Bucket([]byte("all_dates"))
	publicDates := s.Bucket([]byte("public_dates"))

	if err = data.Put([]byte(searchID), bytes); err != nil {
		return err
	}
	if err = allDates.Put([]byte(created), []byte(searchID)); err != nil {
		return err
	}
	if !private {
		if err = publicDates.Put([]byte(created), []byte(searchID)); err != nil {
			return err
		}
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetSearch get Search data by ID
func GetSearch(searchID string) ([]byte, error) {
	var data []byte
	var err error
	err = db.View(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("searches"))
		sd := s.Bucket([]byte("search_data"))
		data = sd.Get([]byte(searchID))
		return nil
	})
	if len(data) == 0 {
		err = errors.New("No data found")
	}
	return data, err
}

// SaveSummary saves the Search Summary to DB
func SaveSummary(searchID string, bytes []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		// Get root Searches bucket
		s := tx.Bucket([]byte("searches"))
		// Get Search Data Bucket
		data := s.Bucket([]byte("search_data"))

		return data.Put([]byte(searchID+"_summary"), bytes)
	})
	return err
}

// GetSummary get Search Summary by ID
func GetSummary(searchID string) ([]byte, error) {
	var data []byte
	var err error
	err = db.Update(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("searches"))
		sd := s.Bucket([]byte("search_data"))
		data = sd.Get([]byte(searchID + "_summary"))
		return nil
	})
	if len(data) == 0 {
		err = errors.New("No data found")
	}
	return data, err
}

// SaveMatches saves the Search Matches to DB
func SaveMatches(searchID string, list map[string][]byte) error {
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get root Searches bucket
	s := tx.Bucket([]byte("searches"))
	// Get Search Data Bucket
	data := s.Bucket([]byte("search_data"))

	for slug, bytes := range list {
		err := data.Put([]byte(searchID+"_matches_"+slug), bytes)
		if err != nil {
			return err
		}
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetMatches get Search Matches data by ID
func GetMatches(searchID string, slug string) ([]byte, error) {
	var data []byte
	var err error
	err = db.Update(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("searches"))
		sd := s.Bucket([]byte("search_data"))
		data = sd.Get([]byte(searchID + "_matches_" + slug))
		return nil
	})
	if len(data) == 0 {
		err = errors.New("No data found")
	}
	return data, err
}
