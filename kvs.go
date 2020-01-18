// Package kvfs provides a simple persistent key-value store.
//
// The API is very simple, you can:
//
// - Put() entries
// - Get() entries
// - Delete() entries
// - Keys() dump all keys in a bucket.
//
// kvfs uses BoltDB for storage.
package kvs

import (
	"errors"
	"time"

	"go.etcd.io/bbolt"
)

// KvS is the key value store. Use the Open() method to create
// one, and Close() it when done.
type KvS struct {
	db *bbolt.DB
}

var (
	// ErrNotFound is returned when the key supplied to a Get or Delete
	// method does not exist in the database.
	ErrNotFound = errors.New("kvs: key not found")

	// ErrBadValue is returned when the value supplied to the Put method
	// is nil.
	ErrBadValue = errors.New("kvs: bad value")

	// ErrBucketNotFound is returned when the bucket name supplied does not exists
	ErrBucketNotFound = errors.New("kvs: bucket not found")
)

// Open a key-value store.
// "filename" is the full path to the database file, any leading directories
// must have been created already. File is created with mode 0640 if needed.
//
// "bucket" is a collections of key/value pairs within the database.
// All keys in a bucket must be unique.
func Open(filename string) (*KvS, error) {
	opts := &bbolt.Options{
		Timeout: 50 * time.Millisecond,
	}

	db, err := bbolt.Open(filename, 0640, opts)
	if err != nil {
		return nil, err
	}

	return &KvS{db: db}, nil
}

// Put an entry into the store.
// The key can be an empty string, but the value cannot be nil - if it is,
// Put() returns ErrBadValue.
func (rcv *KvS) Put(bucket, key string, value []byte) error {
	if value == nil {
		return ErrBadValue
	}

	return rcv.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})

	/*
		return rcv.db.Update(func(tx *bbolt.Tx) error {
			return tx.Bucket([]byte(bucket)).Put([]byte(key), value)
		})*/
}

// Buckets returns a a list of buckets.
func (rcv *KvS) Buckets() []string {
	var res []string
	rcv.db.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bbolt.Bucket) error {
			res = append(res, string(name[:]))
			return nil
		})
	})

	return res
}

// Get an entry from the store.
// If the key is not present in the store, Get returns ErrNotFound.
func (rcv *KvS) Get(bucket, key string) ([]byte, error) {
	var v []byte
	rcv.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotFound
		}

		v = b.Get([]byte([]byte(key)))
		return nil

		/*
			c := tx.Bucket([]byte(rcv.bucket)).Cursor()
			k, v = c.Seek([]byte(key))
			if k == nil || string(k) != key {
				return ErrNotFound
			}*/
	})

	return v, nil
}

// Delete the entry with the given key.
// If no such key is present in the store, it returns ErrNotFound.
func (rcv *KvS) Delete(bucket, key string) error {
	return rcv.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotFound
		}

		c := b.Cursor()
		if k, _ := c.Seek([]byte(key)); k == nil || string(k) != key {
			return ErrNotFound
		}

		return c.Delete()
	})
}

// DeleteBucket deletes a bucket.
// Returns an error if the bucket cannot be found or if the key represents a non-bucket value.
func (rcv *KvS) DeleteBucket(bucket string) error {
	return rcv.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
}

// Keys returns a list of all the keys in the specified bucket.
func (rcv *KvS) Keys(bucket string) []string {
	var res []string
	rcv.db.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
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

// Close closes the key-value store file.
func (rcv *KvS) Close() error {
	return rcv.db.Close()
}

func (rcv *KvS) createBucket(bucket string) error {
	return rcv.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
}
