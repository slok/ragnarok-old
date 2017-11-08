package serializer

import (
	"fmt"

	"github.com/ghodss/yaml"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/log"
)

// YAMLSerializer knows how to serialize objects back and forth using YAML style.
type YAMLSerializer struct {
	asseter    TypeAsserter
	factory    Factory
	discoverer TypeDiscoverer
	typer      Typer
	logger     log.Logger
}

// NewYAMLSerializer returns a new YAMLSerializer object.
func NewYAMLSerializer(typer Typer, factory Factory, logger log.Logger) *YAMLSerializer {
	return &YAMLSerializer{
		asseter:    SafeTypeAsserter,
		factory:    factory,
		typer:      typer,
		discoverer: YAMLTypeDiscoverer,
		logger:     logger,
	}
}

// Encode will encode in YAML the received object in the out argument (writer interface).
// Satisfies Serializer interface.
func (y *YAMLSerializer) Encode(obj api.Object, out interface{}) error {
	w, err := y.asseter.Writer(out)
	if err != nil {
		return err
	}

	// Ensure the object has the correct type.
	if err := y.typer.SetType(obj); err != nil {
		return err
	}
	marshalled, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err := w.Write(marshalled); err != nil {
		return err
	}

	return nil
}

// Decode will decode received data YAML ([]byte type) into an Object.
// Satisfies Serializer interface.
func (y *YAMLSerializer) Decode(data interface{}) (api.Object, error) {
	// Get correct data type.
	bdata, err := y.asseter.ByteArray(data)
	if err != nil {
		return nil, err
	}

	// Decode the object as an object kind  to know what kind of object we need to return.
	tm, err := y.discoverer.Discover(bdata)
	if err != nil {
		return nil, fmt.Errorf("unknown type of object: %s", err)
	}

	// Create the specific object.
	obj, err := y.factory.NewPlainObject(tm)
	if err != nil {
		return nil, err
	}

	// Decode the final object correctly.
	if err := yaml.Unmarshal(bdata, obj); err != nil {
		return nil, err
	}

	return obj, nil
}

// yamlTypeDiscoverer implements the TypeDiscoverer interface for the yaml format.
type yamlTypeDiscoverer struct {
	asserter TypeAsserter
}

// YAMLTypeDiscoverer is a discoverery of object kinds based on the yaml format.
var YAMLTypeDiscoverer = &yamlTypeDiscoverer{
	asserter: SafeTypeAsserter,
}

func (y *yamlTypeDiscoverer) Discover(data interface{}) (api.TypeMeta, error) {
	obj := api.TypeMeta{}
	var b []byte

	b, err := y.asserter.ByteArray(data)
	if err != nil {
		return obj, err
	}

	if err := yaml.Unmarshal(b, &obj); err != nil {
		return obj, err
	}

	if obj.Kind == "" || obj.Version == "" {
		return obj, fmt.Errorf("object kind could not be discoved")
	}

	return obj, nil
}
