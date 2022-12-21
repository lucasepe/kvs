package store

import (
	"errors"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Options are the options for the bbolt store.
type Options struct {
	// Bucket name for storing the key-value pairs.
	BucketName string
	// Path of the DB file.
	Path string

	Timeout time.Duration
}

// New creates a new bbolt store.
// You must call the Close() method on the store when you're done working with it.
func New(options Options) (*Store, error) {
	result := &Store{}

	if options.Path == "" {
		return nil, fmt.Errorf("missed store path")
	}
	if options.Timeout == 0 {
		options.Timeout = 50 * time.Millisecond
	}

	// Open DB
	db, err := bolt.Open(options.Path, 0600, &bolt.Options{
		Timeout: options.Timeout,
	})
	if err != nil {
		return result, err
	}

	result.db = db
	result.bucketName = options.BucketName

	return result, nil
}

var (
	// ErrBucketNotFound is returned when the bucket name supplied does not exists
	ErrBucketNotFound = errors.New("kvs: bucket not found")
)

type Store struct {
	db         *bolt.DB
	bucketName string
}

// Set stores the given value for the given key.
// The key must not be "" and the value must not be nil.
func (s *Store) Set(k string, v []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(s.bucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(k), v)
	})
}

// Get retrieves the stored value for the given key.
// The key must not be "" and the pointer must not be nil.
func (s *Store) Get(k string) (v []byte, err error) {
	var data []byte
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		if b == nil {
			return ErrBucketNotFound
		}

		txData := b.Get([]byte(k))
		// txData is only valid during the transaction.
		// Its value must be copied to make it valid outside of the tx.
		// TODO: Benchmark if it's faster to copy + close tx,
		// or to keep the tx open until unmarshalling is done.
		if txData != nil {
			// `data = append([]byte{}, txData...)` would also work, but the following is more explicit
			data = make([]byte, len(txData))
			copy(data, txData)
		}
		return nil
	})

	return data, err
}

// Delete deletes the stored value for the given key.
// Deleting a non-existing key-value pair does NOT lead to an error.
// The key must not be "".
func (s *Store) Delete(k string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		if b == nil {
			return ErrBucketNotFound
		}
		return b.Delete([]byte(k))
	})
}

// DeleteBucket deletes a bucket.
// Returns an error if the bucket cannot be found or if the key represents a non-bucket value.
func (s *Store) DeleteBucket(bucket string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
}

// Keys returns a list of all the keys in the specified bucket.
func (s *Store) Keys() []string {
	var res []string
	s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(s.bucketName))
		if b == nil {
			return ErrBucketNotFound
		}

		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			//binary.Write(w, binary.LittleEndian, k)
			res = append(res, string(k[:]))
		}

		return nil
	})

	return res
}

// Buckets returns a a list of buckets.
func (s *Store) Buckets() []string {
	var res []string
	s.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			// Count only if the bucket has keys
			if b.Stats().KeyN > 0 {
				res = append(res, string(name[:]))
			}

			return nil
		})
	})

	return res
}

// Close closes the store.
// It must be called to make sure that all open transactions finish and to release all DB resources.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}
