package kvstore

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"path/filepath"
	"runtime"
	"vastchain.ltd/vastchain/chain_structure"
)

// DB is the key-value database of Vastchain.
type DB struct {
	db          *badger.DB
	initialized bool
}

// Tx represents a transaction of Vastchain.
type Tx struct {
	db       *DB
	txn      *badger.Txn
	isUpdate bool
}

var ErrKeyNotFound = chain_structure.NewVcErrorInvalidArgument("key", "key not found")

// NewDB creates a new DB instance.
// If related database is not found, it will be created.
func NewDB(dataPath string) (*DB, error) {
	if dataPath == "" {
		return nil, chain_structure.NewVcErrorInvalidArgument("dataPath", "")
	}
	options := badger.DefaultOptions(filepath.Join(dataPath, "vastchain.db"))
	if runtime.GOOS == "windows" {
		options.WithTruncate(true)
	}
	db, err := badger.Open(options)
	if err != nil {
		return nil, errors.Wrap(err, "KVStore:NewDB: failed to open the database")
	}

	return &DB{db: db, initialized: true}, nil
}

// Dispose disposes all related resources.
// After this call, any call to the DB struct is not allowed.
func (db *DB) Dispose() {
	if !db.initialized {
		return
	}
	db.db.Close()
}

// View opens an read-only transaction.
func (db *DB) View(fn func(tx *Tx) error) error {
	if !db.initialized {
		return chain_structure.NewVcErrorNotInitialized()
	}
	return db.db.View(func(txn *badger.Txn) error {
		tx := &Tx{db: db, txn: txn, isUpdate: false}
		return fn(tx)
	})
}

// Update opens an read-write transaction.
func (db *DB) Update(fn func(tx *Tx) error) error {
	if !db.initialized {
		return chain_structure.NewVcErrorNotInitialized()
	}
	return db.db.Update(func(txn *badger.Txn) error {
		tx := &Tx{db: db, txn: txn, isUpdate: true}
		return fn(tx)
	})
}

// Creates an raw transaction.
func (db *DB) NewTransaction(isUpdate bool) *Tx {
	if !db.initialized {
		return nil
	}
	txn := db.db.NewTransaction(isUpdate)
	return &Tx{
		db:       db,
		txn:      txn,
		isUpdate: isUpdate,
	}
}

// Set sets the value of a key.
func (tx *Tx) Set(bucket string, key []byte, value []byte) error {
	if tx.db == nil || !tx.db.initialized {
		return chain_structure.NewVcErrorNotInitialized()
	}
	if err := checkKey(bucket, key); err != nil {
		return err
	}
	if !tx.isUpdate {
		return chain_structure.NewVcErrorPreconditionNotSatisfied("read-only transaction does not support write")
	}

	return tx.txn.Set(mergeKey(bucket, key), value)
}

// Delete deletes a key.
func (tx *Tx) Delete(bucket string, key []byte) error {
	if tx.db == nil || !tx.db.initialized {
		return chain_structure.NewVcErrorNotInitialized()
	}
	if err := checkKey(bucket, key); err != nil {
		if err == badger.ErrKeyNotFound {
			return ErrKeyNotFound
		}
		return err
	}
	if !tx.isUpdate {
		return chain_structure.NewVcErrorPreconditionNotSatisfied("read-only transaction does not support write")
	}

	return tx.txn.Delete(mergeKey(bucket, key))
}

// Get gets the value from a key.
func (tx *Tx) Get(bucket string, key []byte) ([]byte, error) {
	if tx.db == nil || !tx.db.initialized {
		return nil, chain_structure.NewVcErrorNotInitialized()
	}
	if err := checkKey(bucket, key); err != nil {
		return nil, err
	}
	item, err := tx.txn.Get(mergeKey(bucket, key))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	// TODO: reuse buffer to improve performance
	v, err := item.ValueCopy(nil)

	if err != nil {
		return nil, err
	}
	return v, nil
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	if tx.db == nil || !tx.db.initialized {
		return chain_structure.NewVcErrorNotInitialized()
	}
	return tx.txn.Commit()
}

// Rollback rollbacks the transaction.
func (tx *Tx) Rollback() {
	tx.txn.Discard()
}

// checkKey checks if the bucket and key is valid.
func checkKey(bucket string, key []byte) error {
	lengthBucket := len(bucket)
	if lengthBucket == 0 || lengthBucket > 255 {
		return chain_structure.NewVcErrorInvalidArgument("bucket", "the length should be greater than 0 less than 256")
	}

	lengthKey := len(key)
	if lengthKey == 0 || lengthKey > 255 {
		return chain_structure.NewVcErrorInvalidArgument("key", "the length should be greater than 0 less than 256")
	}

	return nil
}

// mergeKey generates an internal key which combines bucket and key.
// A length prefix is added before bucket value.
func mergeKey(bucket string, key []byte) []byte {
	ret := make([]byte, len(bucket)+len(key)+1)
	ret[0] = byte(len(bucket))
	copy(ret[1:], bucket)
	copy(ret[len(bucket)+1:], key)
	return ret
}
