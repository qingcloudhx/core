package store

import (
	"errors"

	"github.com/boltdb/bolt"
)

const (
	// Permissions to use on the db file. This is only used if the
	// database file does not exist and needs to be created.
	dbFileMode = 0600
)

var (
	// Bucket names we perform transactions in
	//dbLogs = []byte("logs")
	dbConf = []byte("conf")

	// An error indicating a given key does not exist
	ErrKeyNotFound = errors.New("not found")
)

// BoltStore provides access to BoltDB for Raft to store and retrieve
// log entries. It also provides key/value storage, and can be used as
// a LogStore and StableStore.
type BoltStore struct {
	// conn is the underlying handle to the db.
	conn *bolt.DB

	// The path to the Bolt database file
	path string
}

// Options contains all the configuraiton used to open the BoltDB
type Options struct {
	// Path is the file path to the BoltDB to use
	Path string

	// BoltOptions contains any specific BoltDB options you might
	// want to specify [e.g. open timeout]
	BoltOptions *bolt.Options

	// NoSync causes the database to skip fsync calls after each
	// write to the log. This is unsafe, so it should be used
	// with caution.
	NoSync bool
}

// readOnly returns true if the contained bolt options say to open
// the DB in readOnly mode [this can be useful to tools that want
// to examine the log]
func (o *Options) readOnly() bool {
	return o != nil && o.BoltOptions != nil && o.BoltOptions.ReadOnly
}

// NewBoltStore takes a file path and returns a connected Raft backend.
func NewBoltStore(path string) (*BoltStore, error) {
	return New(Options{Path: path})
}

// New uses the supplied options to open the BoltDB and prepare it for use as a raft backend.
func New(options Options) (*BoltStore, error) {
	// Try to connect
	handle, err := bolt.Open(options.Path, dbFileMode, options.BoltOptions)
	if err != nil {
		return nil, err
	}
	handle.NoSync = options.NoSync

	// Create the new store
	store := &BoltStore{
		conn: handle,
		path: options.Path,
	}

	// If the store was opened read-only, don't try and create buckets
	if !options.readOnly() {
		// Set up our buckets
		if err := store.initialize(); err != nil {
			store.Close()
			return nil, err
		}
	}
	return store, nil
}

// initialize is used to set up all of the buckets.
func (b *BoltStore) initialize() error {
	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create all the buckets
	if _, err := tx.CreateBucketIfNotExists(dbConf); err != nil {
		return err
	}

	return tx.Commit()
}

// Close is used to gracefully close the DB connection.
func (b *BoltStore) Close() error {
	return b.conn.Close()
}

func (b *BoltStore) Del(k []byte) error {
	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbConf)
	if err := bucket.Delete(k); err != nil {
		return err
	}

	return tx.Commit()
}

// Set is used to set a key/value set outside of the raft log
func (b *BoltStore) Set(k, v []byte) error {
	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbConf)
	if err := bucket.Put(k, v); err != nil {
		return err
	}

	return tx.Commit()
}
func (b *BoltStore) ForEach(fn func(k, v []byte) error) error {
	tx, err := b.conn.Begin(false)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbConf)
	return bucket.ForEach(fn)
}

// Get is used to retrieve a value from the k/v store by key
func (b *BoltStore) Get(k []byte) ([]byte, error) {
	tx, err := b.conn.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbConf)
	val := bucket.Get(k)

	if val == nil {
		return nil, ErrKeyNotFound
	}
	return append([]byte(nil), val...), nil
}

// SetUint64 is like Set, but handles uint64 values
//func (b *BoltStore) SetUint64(key []byte, val uint64) error {
//	return b.Set(key, uint64ToBytes(val))
//}

// GetUint64 is like Get, but handles uint64 values
//func (b *BoltStore) GetUint64(key []byte) (uint64, error) {
//	val, err := b.Get(key)
//	if err != nil {
//		return 0, err
//	}
//	return bytesToUint64(val), nil
//}

// Sync performs an fsync on the database file handle. This is not necessary
// under normal operation unless NoSync is enabled, in which this forces the
// database file to sync against the disk.
func (b *BoltStore) Sync() error {
	return b.conn.Sync()
}
