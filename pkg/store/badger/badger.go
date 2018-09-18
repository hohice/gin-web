package badger

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/hohice/gin-web/pkg/setting"
	"github.com/hohice/gin-web/pkg/store"
	"github.com/hohice/gin-web/pkg/util/logger"
)

const NAME = "badger"

var (
	log *logger.Logger
)

type BadgerStore struct {
	bdb *badger.DB
}

var badgerStire *BadgerStore

func init() {
	configChan := make(chan struct{})
	setting.RegNotifyChannel(configChan)
	go func() {
		for {
			select {
			case _, ok := <-configChan:
				{
					if !ok {
						return
					} else {
						if conf, ok := setting.Config.Store.Bases[NAME]; ok {
							if conf.Enable {
								log = logger.DefaultLogger
								opts := badger.DefaultOptions
								opts.Dir = conf.BasePath
								opts.ValueDir = conf.ValuePath
								//open db
								db, err := badger.Open(opts)
								if err != nil {
									log.Fatalw(fmt.Sprintf("%v", err), "store", "badger")
									return
								}
								badgerStire = &BadgerStore{
									bdb: db,
								}
								store.SetMethodInstance(badgerStire)
							}
						}
					}
				}
			}
		}
	}()
}

func (bs *BadgerStore) Close() {
	defer bs.bdb.Close()
}

func (bs *BadgerStore) Get(key string) ([]byte, error) {
	var ritem *badger.Item
	err := bs.bdb.View(func(txn *badger.Txn) error {
		if item, err := txn.Get([]byte(key)); err != nil {
			return err
		} else {
			ritem = item
			return nil
		}
	})
	if err != nil {
		return []byte{}, err
	}
	return ritem.Value()
}
func (bs *BadgerStore) Set(key string, value []byte) error {
	err := bs.bdb.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(key), value); err != nil {
			return err
		} else {
			return txn.Commit(func(error) {})
		}
	})
	return err
}
func (bs *BadgerStore) Del(key string) error {
	err := bs.bdb.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(key)); err != nil {
			return err
		} else {
			return txn.Commit(func(error) {})
		}
	})
	return err
}
