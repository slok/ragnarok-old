package apimachinery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ghodss/yaml"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/log"
)

// Serializer has ability to serialize objects back and forward in different formats.
type Serializer interface {
	// Encode will take an object and will write the serialized data on the received writer.
	Encode(obj api.Object, w io.Writer) error
	// Decode will take raw data in format and return a runtime object.
	Decode(data []byte) (api.Object, error)
	// Decode will take raw data in format and return a runtime object.
	DecodeInto(data []byte) (api.Object, error)
}

// JSONSerializer knows how to serialize objects back and forth using JSON style.
type JSONSerializer struct {
	factory    Factory
	discoverer TypeDiscoverer
	typer      Typer
	logger     log.Logger
}

// NewJSONSerializer returns a new JSONSerializer object
func NewJSONSerializer(typer Typer, factory Factory, logger log.Logger) *JSONSerializer {
	return &JSONSerializer{
		factory:    factory,
		typer:      typer,
		discoverer: JSONTypeDiscoverer,
		logger:     logger,
	}
}

// Encode satisfies Serializer interface.
func (j *JSONSerializer) Encode(obj api.Object, w io.Writer) error {
	// Ensure the object has the correct type.
	if err := j.typer.SetType(obj); err != nil {
		return err
	}
	e := json.NewEncoder(w)
	return e.Encode(obj)
}

// Decode satisfies Serializer interface.
func (j *JSONSerializer) Decode(data []byte) (api.Object, error) {
	// Decode the object as an object kind  to know what kind of object we need to return.
	tm, err := j.discoverer.Discover(data)
	if err != nil {
		return nil, fmt.Errorf("unknown type of object: %s", err)
	}

	// Create the specific object.
	obj, err := j.factory.NewPlainObject(tm)
	if err != nil {
		return nil, err
	}

	// Decode the final object correctly.
	d := json.NewDecoder(bytes.NewReader(data))
	if err := d.Decode(obj); err != nil {
		return nil, err
	}

	return obj, nil
}

// YAMLSerializer knows how to serialize objects back and forth using YAML style.
type YAMLSerializer struct {
	factory    Factory
	discoverer TypeDiscoverer
	typer      Typer
	logger     log.Logger
}

// NewYAMLSerializer returns a new YAMLSerializer object
func NewYAMLSerializer(typer Typer, factory Factory, logger log.Logger) *YAMLSerializer {
	return &YAMLSerializer{
		factory:    factory,
		typer:      typer,
		discoverer: YAMLTypeDiscoverer,
		logger:     logger,
	}
}

// Encode satisfies Serializer interface.
func (y *YAMLSerializer) Encode(obj api.Object, w io.Writer) error {
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

// Decode satisfies Serializer interface.
func (y *YAMLSerializer) Decode(data []byte) (api.Object, error) {
	// Decode the object as an object kind  to know what kind of object we need to return.
	tm, err := y.discoverer.Discover(data)
	if err != nil {
		return nil, fmt.Errorf("unknown type of object: %s", err)
	}

	// Create the specific object.
	obj, err := y.factory.NewPlainObject(tm)
	if err != nil {
		return nil, err
	}

	// Decode the final object correctly.
	if err := yaml.Unmarshal(data, obj); err != nil {
		return nil, err
	}

	return obj, nil
}
