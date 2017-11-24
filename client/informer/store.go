package informer

import (
	"sync"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/log"
)

// ObjectIndexKeyer knows how to return an index key based on a resource.
type ObjectIndexKeyer interface {
	// GetKey receives an object and returns an index key.
	GetKey(obj api.Object) (string, error)
}

// ObjectIndexKeyerFunc is a helper so a function can use as an ObjectIndexKeyer.
type ObjectIndexKeyerFunc func(obj api.Object) (string, error)

// GetKey satisfies ObjectIndexKeyer.
func (o ObjectIndexKeyerFunc) GetKey(obj api.Object) (string, error) {
	return o(obj)
}

// Store stores resources in an indexed style.
type Store interface {
	// Add stores an object on the store.
	Add(obj api.Object) error
	// Update stores an object on the store (An alias of add, if the resource doesn't exist id doesn't return an error).
	Update(obj api.Object) error
	// GetByKey returns the resource associated with the key, it also returns if the object exists.
	GetByKey(key string) (obj api.Object, exists bool, err error)
	// Get returns the resource from the resource passed (the resource that is passed will be used for getting the key),
	// it also returns if the object exists.
	Get(obj api.Object) (newObj api.Object, exists bool, err error)
	// Delete deletes the resource.
	Delete(obj api.Object) error
}

// IndexedStore just indexes the resources based on the object using an internal map style. The index
// key will be calculater with the indexer passed to the object. Is safe to use it concurrently but is
// not the fastest.
type IndexedStore struct {
	reg     *sync.Map
	indexer ObjectIndexKeyer
	logger  log.Logger
}

// NewIndexedStore returns a new IndexedStore.
func NewIndexedStore(indexer ObjectIndexKeyer, registry *sync.Map, logger log.Logger) *IndexedStore {
	return &IndexedStore{
		reg:     registry,
		indexer: indexer,
	}
}

// Add satisfies Store interface.
func (i *IndexedStore) Add(obj api.Object) error {
	key, err := i.indexer.GetKey(obj)
	if err != nil {
		return err
	}
	i.reg.Store(key, obj)
	return nil
}

// Update satisfies Store interface.
func (i *IndexedStore) Update(obj api.Object) error {
	return i.Add(obj)
}

// Get satisfies Store interface.
func (i *IndexedStore) Get(obj api.Object) (api.Object, bool, error) {
	key, err := i.indexer.GetKey(obj)
	if err != nil {
		return nil, false, err
	}
	return i.GetByKey(key)
}

// GetByKey satisfies Store interface.
func (i *IndexedStore) GetByKey(key string) (api.Object, bool, error) {
	obj, ok := i.reg.Load(key)
	if !ok {
		return nil, ok, nil
	}
	return obj.(api.Object), ok, nil
}

// Delete satisfies Store interface.
func (i *IndexedStore) Delete(obj api.Object) error {
	key, err := i.indexer.GetKey(obj)
	if err != nil {
		return err
	}
	i.reg.Delete(key)
	return nil
}
