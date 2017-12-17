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
