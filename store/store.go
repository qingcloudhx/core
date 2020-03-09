package store

// StableStore is used to provide stable storage
// of key configurations to ensure safety.
type StableStore interface {
	//init
	ForEach(fn func(k, v []byte) error) error
	Set(key []byte, val []byte) error
	// Get returns the value for key, or an empty byte slice if key was not found.
	Get(key []byte) ([]byte, error)
	//delete
	Del(key []byte) error
	Sync() error
}
