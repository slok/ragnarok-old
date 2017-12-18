package api

import (
	"github.com/slok/ragnarok/api"
)

type TestObj struct {
	Kind    api.Kind
	Version api.Version
	ID      string
	Labels  map[string]string
}

func (t *TestObj) GetObjectKind() api.Kind       { return t.Kind }
func (t *TestObj) GetObjectVersion() api.Version { return t.Version }
func (t *TestObj) GetObjectMetadata() api.ObjectMeta {
	return api.ObjectMeta{ID: t.ID, Labels: t.Labels}
}
func (t *TestObj) DeepCopy() api.Object {
	copy := *t
	return &copy
}

type TestObjList struct {
	Kind     api.Kind
	Version  api.Version
	Continue string
	Items    []*TestObj
}

func NewTestObjList(testObjs []*TestObj, continueList string) TestObjList {
	return TestObjList{
		Continue: continueList,
		Items:    testObjs,
	}
}

func (t *TestObjList) GetObjectKind() api.Kind       { return t.Kind }
func (t *TestObjList) GetObjectVersion() api.Version { return t.Version }

// GetObjectMetadata satisfies object interface.
func (t *TestObjList) GetObjectMetadata() api.ObjectMeta {
	return api.NoObjectMeta
}

// GetListMetadata satisfies objectList interface.
func (t *TestObjList) GetListMetadata() api.ListMeta {
	return api.ListMeta{Continue: t.Continue}
}

// GetItems satisfies ObjectList interface.
func (t *TestObjList) GetItems() []api.Object {
	res := make([]api.Object, len(t.Items))
	for i, item := range t.Items {
		res[i] = api.Object(item)
	}
	return res
}

// DeepCopy satisfies object interface.
func (t *TestObjList) DeepCopy() api.Object {
	ts := []*TestObj{}
	for i, testObj := range t.Items {
		t := *testObj
		ts[i] = &t
	}
	copy := NewTestObjList(ts, t.Continue)
	return &copy
}
