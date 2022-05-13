package main

import "github.com/boltdb/bolt"

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

func (storage BoltStorage) StoreSquidForuserId(userId string, squid string) error {
	return storage.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("squidstate"))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(userId), []byte(squid))
	})
}

func (storage BoltStorage) GetSquidForUserId(userId string) (squid string, err error) {
	return squid, storage.db.View(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("squidstate"))
		if err != nil {
			return err
		}
		squid = string(bucket.Get([]byte(userId)))
		return nil
	})
}

func (storage BoltStorage) IsReactableMessage(messageId string) (ok bool, err error) {
	return ok, storage.db.View(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("askstate"))
		if err != nil {
			return err
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
