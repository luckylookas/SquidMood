package storage

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

type BoltStorage struct {
	db *bolt.DB
}

func New(file string) (BoltStorage, error) {
	db, err := bolt.Open(file, 0600, nil)
	return BoltStorage{db: db}, err
}

func (storage BoltStorage) Close() {
	storage.db.Close()
}

func (storage BoltStorage) StoreSquidForUserId(userId string, squid string) error {
	return storage.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("squidstate"))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(fmt.Sprintf("%s", userId)), []byte(squid))
	})
}

func (storage BoltStorage) GetSquidForUserId(userId string) (squid string, err error) {
	return squid, storage.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("squidstate"))
		if bucket == nil {
			return errors.New("no squid states yet")
		}
		squid = string(bucket.Get([]byte(fmt.Sprintf("%s", userId))))
		return nil
	})
}

func (storage BoltStorage) IsReactableMessage(messageId string) (ok bool, err error) {
	return ok, storage.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("askstate"))
		if bucket == nil {
			return errors.New("no ask states yet")
		}
		ok = string(bucket.Get([]byte(messageId))) != ""
		return nil
	})
}

func (storage BoltStorage) StoreReactableMessage(messageId string) error {
	return storage.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("askstate"))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(messageId), []byte(messageId))
	})
}
