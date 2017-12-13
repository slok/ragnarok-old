package store_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/client/util/store"
	"github.com/slok/ragnarok/log"
	mstore "github.com/slok/ragnarok/mocks/client/util/store"
	testapi "github.com/slok/ragnarok/test/api"
)

func TestIndexedStoreAdd(t *testing.T) {
	tests := []struct {
		name     string
		indexKey string
		obj      api.Object
		keyErr   bool
		expErr   bool
	}{
		{
			name:     "Storing a correct object with a correct key should store the object on the registry.",
			indexKey: "/test/v1/testobj1",
			obj:      &testapi.TestObj{ID: "testobj1"},
			keyErr:   false,
			expErr:   false,
		},
		{
			name:     "Storing a correct object with an incorrect key should return an error.",
			indexKey: "",
			obj:      &testapi.TestObj{ID: "testobj1"},
			keyErr:   true,
			expErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var keyErr error
			if test.keyErr {
				keyErr = errors.New("wanted error")
			}

			// Mocks.
			mi := &mstore.ObjectIndexKeyer{}
			mi.On("GetKey", mock.Anything).Return(test.indexKey, keyErr)

			// Create store.
			reg := &sync.Map{}
			store := store.NewIndexedStore(mi, reg, log.Dummy)

			// Add the required objects
			err := store.Add(test.obj)

			// Check the objects are there.
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				gotObj, ok := reg.Load(test.indexKey)
				assert.True(ok)
				assert.Equal(test.obj, gotObj)
			}
		})
	}
}

func TestIndexedStoreUpdate(t *testing.T) {
	tests := []struct {
		name     string
		indexKey string
		obj      api.Object
		keyErr   bool
		expErr   bool
	}{
		{
			name:     "Storing a correct object with a correct key should store the object on the registry.",
			indexKey: "/test/v1/testobj1",
			obj:      &testapi.TestObj{ID: "testobj1"},
			keyErr:   false,
			expErr:   false,
		},
		{
			name:     "Storing a correct object with an incorrect key should return an error.",
			indexKey: "",
			obj:      &testapi.TestObj{ID: "testobj1"},
			keyErr:   true,
			expErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var keyErr error
			if test.keyErr {
				keyErr = errors.New("wanted error")
			}

			// Mocks.
			mi := &mstore.ObjectIndexKeyer{}
			mi.On("GetKey", mock.Anything).Return(test.indexKey, keyErr)

			// Create store.
			reg := &sync.Map{}
			store := store.NewIndexedStore(mi, reg, log.Dummy)

			// Add the required objects
			err := store.Update(test.obj)

			// Check the objects are there.
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				gotObj, ok := reg.Load(test.indexKey)
				assert.True(ok)
				assert.Equal(test.obj, gotObj)
			}
		})
	}
}

func TestIndexedStoreGet(t *testing.T) {
	tests := []struct {
		name      string
		registry  map[string]api.Object
		indexKey  string
		obj       api.Object
		expObj    api.Object
		expExists bool
		keyErr    bool
		expErr    bool
	}{
		{
			name: "Retrieving a present object should retrieve the resource correctly.",
			registry: map[string]api.Object{
				"/test/v1/testobj1": &testapi.TestObj{ID: "testobj1"},
			},
			indexKey:  "/test/v1/testobj1",
			obj:       &testapi.TestObj{ID: "testobj1"},
			expObj:    &testapi.TestObj{ID: "testobj1"},
			expExists: true,
			keyErr:    false,
			expErr:    false,
		},
		{
			name: "Retrieving a non present object should return that the key doesnt exists but without error.",
			registry: map[string]api.Object{
				"/test/v1/testobj2": &testapi.TestObj{ID: "testobj2"},
			},
			indexKey:  "/test/v1/testobj1",
			obj:       &testapi.TestObj{ID: "testobj1"},
			expObj:    nil,
			expExists: false,
			keyErr:    false,
			expErr:    false,
		},
		{
			name: "Retrieving a present object but with a key error should return an error.",
			registry: map[string]api.Object{
				"/test/v1/testobj1": &testapi.TestObj{ID: "testobj1"},
			},
			indexKey:  "/test/v1/testobj1",
			obj:       &testapi.TestObj{ID: "testobj1"},
			expObj:    nil,
			expExists: false,
			keyErr:    true,
			expErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var keyErr error
			if test.keyErr {
				keyErr = errors.New("wanted error")
			}

			// Mocks.
			mi := &mstore.ObjectIndexKeyer{}
			mi.On("GetKey", mock.Anything).Return(test.indexKey, keyErr)

			// Create store.
			reg := &sync.Map{}
			for k, v := range test.registry {
				reg.Store(k, v)
			}
			store := store.NewIndexedStore(mi, reg, log.Dummy)

			// Add the required objects
			gotObj, ex, err := store.Get(test.obj)

			// Check the objects are there.
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(test.expExists, ex, "exists of the resource failed")
				assert.Equal(test.expObj, gotObj)
			}
		})
	}
}

func TestIndexedStoreDelete(t *testing.T) {
	tests := []struct {
		name        string
		registry    map[string]api.Object
		indexKey    string
		obj         api.Object
		expRegistry map[string]api.Object
		keyErr      bool
		expErr      bool
	}{
		{
			name: "Deleting a present object should delete the object.",
			registry: map[string]api.Object{
				"/test/v1/testobj1": &testapi.TestObj{ID: "testobj1"},
				"/test/v1/testobj2": &testapi.TestObj{ID: "testobj2"},
				"/test/v1/testobj3": &testapi.TestObj{ID: "testobj3"},
			},
			indexKey: "/test/v1/testobj2",
			obj:      &testapi.TestObj{ID: "testobj2"},
			expRegistry: map[string]api.Object{
				"/test/v1/testobj1": &testapi.TestObj{ID: "testobj1"},
				"/test/v1/testobj3": &testapi.TestObj{ID: "testobj3"},
			},
			keyErr: false,
			expErr: false,
		},
		{
			name: "Deleting a not present object should not return an error.",
			registry: map[string]api.Object{
				"/test/v1/testobj1": &testapi.TestObj{ID: "testobj1"},
				"/test/v1/testobj3": &testapi.TestObj{ID: "testobj3"},
			},
			indexKey: "/test/v1/testobj2",
			obj:      &testapi.TestObj{ID: "testobj2"},
			expRegistry: map[string]api.Object{
				"/test/v1/testobj1": &testapi.TestObj{ID: "testobj1"},
				"/test/v1/testobj3": &testapi.TestObj{ID: "testobj3"},
			},
			keyErr: false,
			expErr: false,
		},
		{
			name: "Deleting a present object with invalid key should return an error.",
			registry: map[string]api.Object{
				"/test/v1/testobj1": &testapi.TestObj{ID: "testobj1"},
				"/test/v1/testobj2": &testapi.TestObj{ID: "testobj2"},
				"/test/v1/testobj3": &testapi.TestObj{ID: "testobj3"},
			},
			indexKey: "/test/v1/testobj2",
			obj:      &testapi.TestObj{ID: "testobj2"},
			expRegistry: map[string]api.Object{
				"/test/v1/testobj1": &testapi.TestObj{ID: "testobj1"},
				"/test/v1/testobj3": &testapi.TestObj{ID: "testobj3"},
			},
			keyErr: true,
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var keyErr error
			if test.keyErr {
				keyErr = errors.New("wanted error")
			}

			// Mocks.
			mi := &mstore.ObjectIndexKeyer{}
			mi.On("GetKey", mock.Anything).Return(test.indexKey, keyErr)

			// Create store.
			reg := &sync.Map{}
			for k, v := range test.registry {
				v := *(v.(*testapi.TestObj))
				reg.Store(k, &v)
			}
			store := store.NewIndexedStore(mi, reg, log.Dummy)

			// Add the required objects
			err := store.Delete(test.obj)

			// Check the objects are there.
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				// Create our map representation of the sync map so we can compare with the expected one.
				gotReg := map[string]api.Object{}
				reg.Range(func(k, v interface{}) bool {
					key := k.(string)
					value := v.(api.Object)
					gotReg[key] = value
					return true
				})

				assert.Equal(test.expRegistry, gotReg)
			}
		})
	}
}
