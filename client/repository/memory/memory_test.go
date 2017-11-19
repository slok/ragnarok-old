package memory_test

import (
	"sort"
	"testing"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/client/repository/memory"
	"github.com/slok/ragnarok/log"
	testapi "github.com/slok/ragnarok/test/api"
	"github.com/stretchr/testify/assert"
)

func TestMemoryRepositoryCreate(t *testing.T) {
	tests := []struct {
		name     string
		registry map[string]map[string]api.Object
		obj      *testapi.TestObj
		expErr   bool
	}{
		{
			name:     "One new object should be created without error.",
			registry: map[string]map[string]api.Object{},
			obj:      &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1"},
			expErr:   false,
		},
		{
			name: "Storing a object that already was stored should return an error.",
			registry: map[string]map[string]api.Object{
				"testing/v1/test": {
					"test1": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1"},
				},
			},
			obj:    &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1"},
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			assert := assert.New(t)

			cli := memory.NewClient(test.registry, log.Dummy)
			_, err := cli.Create(test.obj)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				nGot, ok := test.registry["testing/v1/test"][test.obj.ID]
				if assert.True(ok) {
					assert.Equal(test.obj, nGot)
				}
			}
		})
	}
}

func TestMemoryRepositoryUdate(t *testing.T) {
	tests := []struct {
		name        string
		registry    map[string]map[string]api.Object
		obj         *testapi.TestObj
		expErr      bool
		expRegistry map[string]map[string]api.Object
	}{
		{
			name:     "When updating a new object sould return an error.",
			registry: map[string]map[string]api.Object{},
			obj:      &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1"},
			expErr:   true,
		},
		{
			name: "Storing a node that already was stored shouldupdate ok.",
			registry: map[string]map[string]api.Object{
				"testing/v1/test": {
					"test1": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1",
						Labels: map[string]string{"test": "wrong"},
					},
				},
			},
			obj: &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1",
				Labels: map[string]string{"test": "good"},
			},
			expRegistry: map[string]map[string]api.Object{
				"testing/v1/test": {
					"test1": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1",
						Labels: map[string]string{"test": "good"},
					},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			cli := memory.NewClient(test.registry, log.Dummy)
			_, err := cli.Update(test.obj)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(test.expRegistry, test.registry)
			}
		})
	}
}

func TestMemoryRepositoryDelete(t *testing.T) {
	tests := []struct {
		name         string
		registry     map[string]map[string]api.Object
		deleteFullID string
		expErr       bool
		expRegistry  map[string]map[string]api.Object
	}{
		{
			name:         "Deleting a missing object shouldn't return an error.",
			registry:     map[string]map[string]api.Object{},
			deleteFullID: "testing/v1/test/test1",
			expRegistry:  map[string]map[string]api.Object{},
			expErr:       false,
		},
		{
			name: "Deleting a not missing object shouldn't return an error.",
			registry: map[string]map[string]api.Object{
				"testing/v1/test": {
					"test1": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1"},
					"test2": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test2"},
				},
			},
			deleteFullID: "testing/v1/test/test2",
			expRegistry: map[string]map[string]api.Object{
				"testing/v1/test": {
					"test1": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1"},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			cli := memory.NewClient(test.registry, log.Dummy)
			err := cli.Delete(test.deleteFullID)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(test.expRegistry, test.registry)
			}
		})
	}
}

func TestMemoryRepositoryGet(t *testing.T) {
	tests := []struct {
		name      string
		registry  map[string]map[string]api.Object
		getFullID string
		expObj    *testapi.TestObj
		expErr    bool
	}{
		{
			name:      "Getting a non existent object should error.",
			registry:  map[string]map[string]api.Object{},
			getFullID: "testing/v1/test/test1",
			expErr:    true,
		},
		{
			name: "Getting an existent object shouldn't error.",
			registry: map[string]map[string]api.Object{
				"testing/v1/test": {
					"test1": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1"},
				},
			},
			getFullID: "testing/v1/test/test1",
			expObj:    &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1"},
			expErr:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			assert := assert.New(t)

			cli := memory.NewClient(test.registry, log.Dummy)
			obj, err := cli.Get(test.getFullID)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(test.expObj, obj)
			}
		})
	}
}

func TestMemoryRepositoryList(t *testing.T) {
	tests := []struct {
		name     string
		registry map[string]map[string]api.Object
		opts     api.ListOptions
		expObjs  []*testapi.TestObj
		expErr   bool
	}{
		{
			name: "No selectors should return all the registry.",
			registry: map[string]map[string]api.Object{
				"testing/v1/test": {
					"test1": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1",
						Labels: map[string]string{"kind": "test-a", "region": "eu-west-1"},
					},
					"test2": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test2",
						Labels: map[string]string{"kind": "test-b", "region": "eu-central-1"},
					},
					"test3": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test3",
						Labels: map[string]string{"kind": "test-a", "region": "eu-central-1"},
					},
				},
				"testing/v2/test": {
					"test4": &testapi.TestObj{Version: "testing/v2", Kind: "test", ID: "test4"},
				},
			},
			opts: api.ListOptions{TypeMeta: api.TypeMeta{Version: "testing/v1", Kind: "test"}},
			expObjs: []*testapi.TestObj{
				&testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1",
					Labels: map[string]string{"kind": "test-a", "region": "eu-west-1"},
				},
				&testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test2",
					Labels: map[string]string{"kind": "test-b", "region": "eu-central-1"},
				},
				&testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test3",
					Labels: map[string]string{"kind": "test-a", "region": "eu-central-1"},
				},
			},
			expErr: false,
		},
		{
			name: "One selectors should return only the selector ones.",
			registry: map[string]map[string]api.Object{
				"testing/v1/test": {
					"test1": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1",
						Labels: map[string]string{"kind": "test-a", "region": "eu-west-1"},
					},
					"test2": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test2",
						Labels: map[string]string{"kind": "test-b", "region": "eu-central-1"},
					},
					"test3": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test3",
						Labels: map[string]string{"kind": "test-a", "region": "eu-central-1"},
					},
				},
				"testing/v2/test": {
					"test4": &testapi.TestObj{Version: "testing/v2", Kind: "test", ID: "test4"},
				},
			},
			opts: api.ListOptions{
				TypeMeta:      api.TypeMeta{Version: "testing/v1", Kind: "test"},
				LabelSelector: map[string]string{"region": "eu-central-1"},
			},
			expObjs: []*testapi.TestObj{
				&testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test2",
					Labels: map[string]string{"kind": "test-b", "region": "eu-central-1"},
				},
				&testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test3",
					Labels: map[string]string{"kind": "test-a", "region": "eu-central-1"},
				},
			},
			expErr: false,
		},
		{
			name: "Multiple selectors should return only the selector ones.",
			registry: map[string]map[string]api.Object{
				"testing/v1/test": {
					"test1": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test1",
						Labels: map[string]string{"kind": "test-a", "region": "eu-west-1"},
					},
					"test2": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test2",
						Labels: map[string]string{"kind": "test-b", "region": "eu-central-1"},
					},
					"test3": &testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test3",
						Labels: map[string]string{"kind": "test-a", "region": "eu-central-1"},
					},
				},
				"testing/v2/test": {
					"test4": &testapi.TestObj{Version: "testing/v2", Kind: "test", ID: "test4"},
				},
			},
			opts: api.ListOptions{
				TypeMeta:      api.TypeMeta{Version: "testing/v1", Kind: "test"},
				LabelSelector: map[string]string{"kind": "test-a", "region": "eu-central-1"},
			},
			expObjs: []*testapi.TestObj{
				&testapi.TestObj{Version: "testing/v1", Kind: "test", ID: "test3",
					Labels: map[string]string{"kind": "test-a", "region": "eu-central-1"},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			assert := assert.New(t)

			cli := memory.NewClient(test.registry, log.Dummy)
			objs, err := cli.List(test.opts)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				sort.Slice(objs, func(i, j int) bool {
					return objs[i].GetObjectMetadata().ID < objs[j].GetObjectMetadata().ID
				})
				oss := make([]*testapi.TestObj, len(objs))
				for i, obj := range objs {
					oss[i] = obj.(*testapi.TestObj)
				}
				assert.Equal(test.expObjs, oss)
			}
		})
	}
}
