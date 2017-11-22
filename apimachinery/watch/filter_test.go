package watch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/apimachinery/watch"
	testapi "github.com/slok/ragnarok/test/api"
)

func TestNoFilter(t *testing.T) {
	tests := []struct {
		name string
		obj  api.Object
	}{
		{
			name: "test object should return false",
			obj:  &testapi.TestObj{Kind: "test", Version: "test/v1"},
		},
		{
			name: "test2 object should return false",
			obj:  &testapi.TestObj{Kind: "test", Version: "test2/v1"},
		},
		{
			name: "test3 object should return false",
			obj:  &testapi.TestObj{Kind: "test", Version: "test3/v1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			res := watch.NoFilter.Filter(test.obj)
			assert.False(res)
		})
	}
}

func TestTypeFilter(t *testing.T) {
	tests := []struct {
		name      string
		filterObj api.Object
		obj       api.Object
		exp       bool
	}{
		{
			name:      "If the filter object type and the filtered object argument are the same type it shouldn't filter",
			filterObj: &testapi.TestObj{Kind: "test", Version: "test/v1"},
			obj:       &testapi.TestObj{Kind: "test", Version: "test/v1"},
			exp:       false,
		},
		{
			name:      "If the filter object type and the filtered object argument are different type it should filter",
			filterObj: &testapi.TestObj{Kind: "test", Version: "test/v2"},
			obj:       &testapi.TestObj{Kind: "test", Version: "test/v1"},
			exp:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			f := watch.NewTypeFilter(test.filterObj)
			res := f.Filter(test.obj)
			assert.Equal(test.exp, res)
		})
	}
}

func TestLabelFilter(t *testing.T) {
	tests := []struct {
		name     string
		selector map[string]string
		obj      api.Object
		exp      bool
	}{
		{
			name:     "Using not selector it shouldn't filter anything",
			selector: map[string]string{},
			obj:      &testapi.TestObj{Labels: map[string]string{"app": "test1"}},
			exp:      false,
		},
		{
			name:     "Using the same selector as the object labels it shouldn't filter the object",
			selector: map[string]string{"app": "test1"},
			obj:      &testapi.TestObj{Labels: map[string]string{"app": "test1"}},
			exp:      false,
		},
		{
			name:     "Using different selector as the object labels it should filter the object",
			selector: map[string]string{"app": "test1"},
			obj:      &testapi.TestObj{Labels: map[string]string{"app": "test2"}},
			exp:      true,
		},
		{
			name:     "Using different and some same selector as the object labels it should filter the object",
			selector: map[string]string{"app": "test1", "test-kind": "whatever"},
			obj:      &testapi.TestObj{Labels: map[string]string{"app": "test1"}},
			exp:      true,
		},
		{
			name:     "Using multiple selector as the object labels it shouldn't filter the object",
			selector: map[string]string{"app": "test1", "test-kind": "whatever"},
			obj:      &testapi.TestObj{Labels: map[string]string{"app": "test1", "test-kind": "whatever", "test-kind2": "whatever2"}},
			exp:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			f := watch.NewLabelFilter(test.selector)
			res := f.Filter(test.obj)
			assert.Equal(test.exp, res)
		})
	}
}

func TestListOptionsFilter(t *testing.T) {
	tests := []struct {
		name string
		opts api.ListOptions
		obj  api.Object
		exp  bool
	}{
		{
			name: "Using no selector and same type it shouldn't filter anything",
			opts: api.ListOptions{
				TypeMeta:      api.TypeMeta{Kind: "test", Version: "test/v1"},
				LabelSelector: map[string]string{},
			},
			obj: &testapi.TestObj{Kind: "test", Version: "test/v1", Labels: map[string]string{"app": "test1"}},
			exp: false,
		},
		{
			name: "Using no selector and different type it should filter",
			opts: api.ListOptions{
				TypeMeta:      api.TypeMeta{Kind: "test2", Version: "test/v1"},
				LabelSelector: map[string]string{},
			},
			obj: &testapi.TestObj{Kind: "test", Version: "test/v1", Labels: map[string]string{"app": "test1"}},
			exp: true,
		},
		{
			name: "Using same selector and same type it shouldn't filter",
			opts: api.ListOptions{
				TypeMeta:      api.TypeMeta{Kind: "test", Version: "test/v1"},
				LabelSelector: map[string]string{"app": "test1"},
			},
			obj: &testapi.TestObj{Kind: "test", Version: "test/v1", Labels: map[string]string{"app": "test1"}},
			exp: false,
		},
		{
			name: "Using same selector and different type it should filter",
			opts: api.ListOptions{
				TypeMeta:      api.TypeMeta{Kind: "test2", Version: "test/v1"},
				LabelSelector: map[string]string{"app": "test1"},
			},
			obj: &testapi.TestObj{Kind: "test", Version: "test/v1", Labels: map[string]string{"app": "test1"}},
			exp: true,
		},
		{
			name: "Using not valid multiple selector and same type it should filter",
			opts: api.ListOptions{
				TypeMeta:      api.TypeMeta{Kind: "test", Version: "test/v1"},
				LabelSelector: map[string]string{"app": "test1", "test-kind": "whatever"},
			},
			obj: &testapi.TestObj{Kind: "test", Version: "test/v1", Labels: map[string]string{"app": "test1"}},
			exp: true,
		},
		{
			name: "Using valid multiple selector and same type it shouldn't filter",
			opts: api.ListOptions{
				TypeMeta:      api.TypeMeta{Kind: "test", Version: "test/v1"},
				LabelSelector: map[string]string{"app": "test1"},
			},
			obj: &testapi.TestObj{Kind: "test", Version: "test/v1", Labels: map[string]string{"app": "test1", "test-kind": "whatever"}},
			exp: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			f := watch.NewListOptionsFilter(test.opts)
			res := f.Filter(test.obj)
			assert.Equal(test.exp, res)
		})
	}
}
