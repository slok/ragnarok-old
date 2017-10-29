package apimachinery

import (
	"encoding/json"
	"fmt"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
)

// Factory implements a way of obtaining objects based on the verison and kind.
type Factory interface {
	// NewPlainObject returns anew plain object based on the version and kind.
	NewPlainObject(t api.TypeMeta) (api.Object, error)
}

// objFactory will be used to get plain objects based on the version and object kind
type objFactory struct{}

// ObjFactory is a global helper to use it as a factory in the econders/decoders.
var ObjFactory = &objFactory{}

// NewPlainObject satisfies Factory interface.
func (o *objFactory) NewPlainObject(t api.TypeMeta) (api.Object, error) {
	// TODO: Make more elegant way of registering object creators.
	switch {
	case t.Kind == clusterv1.NodeKind && t.Version == clusterv1.NodeVersion:
		n := clusterv1.NewNode()
		return &n, nil
	default:
		return nil, fmt.Errorf("unknown %s object type", t)
	}
}

// TypeDiscoverer discovers the type of an object from the encoded format.
type TypeDiscoverer interface {
	// Discovery will return the type of the object.
	Discover(b []byte) (api.TypeMeta, error)
}

// jsonTypeDiscoverer implements the TypeDiscoverer interface for the json format.
type jsonTypeDiscoverer struct{}

// JSONTypeDiscoverer is a discoverery of object kinds based on the json format.
var JSONTypeDiscoverer = &jsonTypeDiscoverer{}

func (j *jsonTypeDiscoverer) Discover(b []byte) (api.TypeMeta, error) {
	obj := api.TypeMeta{}

	if err := json.Unmarshal(b, &obj); err != nil {
		return obj, err
	}

	if obj.Kind == "" || obj.Version == "" {
		return obj, fmt.Errorf("object kind could not be discoved")
	}

	return obj, nil
}

// Typer is the interface that knows to set the type in an object instance.
type Typer interface {
	// SetType sets the type of the object in the object.
	SetType(obj api.Object) error
}

// objTyper is the default object typer.
type objTyper struct{}

// ObjTyper is a handy instance of the default object typer.
var ObjTyper = &objTyper{}

// SetType implements Typer interface.
func (o *objTyper) SetType(obj api.Object) error {
	// TODO: Make more elegant way of setting correct types.
	switch v := obj.(type) {
	case *clusterv1.Node:
		v.Kind = clusterv1.NodeKind
		v.Version = clusterv1.NodeVersion
	default:
		return fmt.Errorf("could not set the type of object because isn't a valid object type")
	}

	return nil
}
