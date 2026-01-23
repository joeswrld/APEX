package storage

import (
	"github.com/dgraph-io/badger/v4"
)

// Database represents the main database interface
type Database struct {
	db *badger.DB
}

// NewDatabase creates a new database instance
func NewDatabase(path string) (*Database, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // Disable logging for production
	
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	
	return &Database{db: db}, nil
}

// Close closes the database
func (d *Database) Close() error {
	return d.db.Close()
}

// Put stores a key-value pair
func (d *Database) Put(key, value []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Get retrieves a value by key
func (d *Database) Get(key []byte) ([]byte, error) {
	var value []byte
	
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		
		value, err = item.ValueCopy(nil)
		return err
	})
	
	return value, err
}

// Delete removes a key-value pair
func (d *Database) Delete(key []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Has checks if a key exists
func (d *Database) Has(key []byte) (bool, error) {
	err := d.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		return err
	})
	
	if err == badger.ErrKeyNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	
	return true, nil
}

// Batch executes multiple operations atomically
func (d *Database) Batch(ops []BatchOp) error {
	return d.db.Update(func(txn *badger.Txn) error {
		for _, op := range ops {
			switch op.Type {
			case BatchOpPut:
				if err := txn.Set(op.Key, op.Value); err != nil {
					return err
				}
			case BatchOpDelete:
				if err := txn.Delete(op.Key); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// Iterator creates an iterator for range queries
func (d *Database) Iterator(prefix []byte) *Iterator {
	txn := d.db.NewTransaction(false)
	opts := badger.DefaultIteratorOptions
	opts.Prefix = prefix
	
	it := txn.NewIterator(opts)
	return &Iterator{
		txn: txn,
		it:  it,
	}
}

// BatchOpType represents batch operation type
type BatchOpType int

const (
	BatchOpPut BatchOpType = iota
	BatchOpDelete
)

// BatchOp represents a batch operation
type BatchOp struct {
	Type  BatchOpType
	Key   []byte
	Value []byte
}

// Iterator wraps badger iterator
type Iterator struct {
	txn *badger.Txn
	it  *badger.Iterator
}

// Next advances iterator
func (i *Iterator) Next() {
	i.it.Next()
}

// Valid checks if iterator is valid
func (i *Iterator) Valid() bool {
	return i.it.Valid()
}

// Key returns current key
func (i *Iterator) Key() []byte {
	return i.it.Item().KeyCopy(nil)
}

// Value returns current value
func (i *Iterator) Value() ([]byte, error) {
	return i.it.Item().ValueCopy(nil)
}

// Close closes iterator
func (i *Iterator) Close() {
	i.it.Close()
	i.txn.Discard()
}

// Rewind moves iterator to start
func (i *Iterator) Rewind() {
	i.it.Rewind()
}