package db

import (
	"github.com/brenordv/go-request/internal/utils"
	"github.com/dgraph-io/badger/v3"
	"os"
	"path"
	"sync"
)

type DatabaseClient struct {
	db *badger.DB
	lock sync.WaitGroup
	WaitForPromises bool
}

func NewDatabaseClient() (*DatabaseClient, error) {
	var appDir string
	var err error
	var db *badger.DB

	appDir, err = utils.GetAppDir()
	if err != nil {
		return nil, err
	}

	dbDir := path.Join(appDir, ".appdata")

	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	opts := badger.DefaultOptions(dbDir)
	opts.CompactL0OnClose = true
	opts.Logger = nil

	db, err = badger.Open(opts)

	return &DatabaseClient{db: db}, err
}

func (dc *DatabaseClient) Add(key []byte, value []byte) error {
	if dc.WaitForPromises {
		defer dc.lock.Done()
	}
	err := dc.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})

	return err
}

func (dc *DatabaseClient) Close() error {
	if dc.WaitForPromises {
		dc.lock.Wait()
	}

	return dc.db.Close()
}

func (dc *DatabaseClient) Get(key []byte) ([]byte, error) {
	var value []byte
	if dc.WaitForPromises {
		dc.lock.Wait()
	}
	err := dc.db.View(func(txn *badger.Txn) error {
		dbRow, err := txn.Get(key)
		if err != nil {
			return err
		}
		return dbRow.Value(func(val []byte) error {
			value = val
			return nil
		})

	})

	return value, err
}

func (dc *DatabaseClient) PromiseWillAdd(delta int) {
	dc.lock.Add(delta)
}