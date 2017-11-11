package serializer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/log"
)

// JSONSerializer knows how to serialize objects back and forth using JSON style.
type JSONSerializer struct {
	asseter    TypeAsserter
	factory    Factory
	discoverer TypeDiscoverer
	typer      Typer
	logger     log.Logger
}

// NewJSONSerializer returns a new JSONSerializer object
func NewJSONSerializer(typer Typer, factory Factory, logger log.Logger) *JSONSerializer {
	return &JSONSerializer{
		asseter:    SafeTypeAsserter,
		factory:    factory,
		typer:      typer,
		discoverer: JSONTypeDiscoverer,
		logger:     logger,
	}
}

// Encode will encode in JSON the received object in the out argument (Writer interface).
// Satisfies Serializer interface.
func (j *JSONSerializer) Encode(obj api.Object, out interface{}) error {
	w, err := j.asseter.Writer(out)
	if err != nil {
		return err
	}

	// Ensure the object has the correct type.
	if err := j.typer.SetType(obj); err != nil {
		return err
	}
	e := json.NewEncoder(w)
	return e.Encode(obj)
}

// Decode will decode received data JSON ([]byte type) into an Object.
// Satisfies Serializer interface.
func (j *JSONSerializer) Decode(data interface{}) (api.Object, error) {
	// Get correct data type.
	bdata, err := j.asseter.ByteArray(data)
	if err != nil {
		return nil, err
	}

	// Decode the object as an object kind  to know what kind of object we need to return.
	tm, err := j.discoverer.Discover(bdata)
	if err != nil {
		return nil, fmt.Errorf("unknown type of object: %s", err)
	}

	// Create the specific object.
	obj, err := j.factory.NewPlainObject(tm)
	if err != nil {
		return nil, err
	}

	// Decode the final object correctly.
	d := json.NewDecoder(bytes.NewReader(bdata))
	if err := d.Decode(obj); err != nil {
		return nil, err
	}

	return obj, nil
}

// jsonTypeDiscoverer implements the TypeDiscoverer interface for the json format.
type jsonTypeDiscoverer struct {
	asserter TypeAsserter
}

// JSONTypeDiscoverer is a discoverery of object kinds based on the json format.
var JSONTypeDiscoverer = &jsonTypeDiscoverer{
	asserter: SafeTypeAsserter,
}

func (j *jsonTypeDiscoverer) Discover(data interface{}) (api.TypeMeta, error) {
	obj := api.TypeMeta{}
	var b []byte

	b, err := j.asserter.ByteArray(data)
	if err != nil {
		return obj, err
	}

	if err := json.Unmarshal(b, &obj); err != nil {
		return obj, err
	}

	if obj.Kind == "" || obj.Version == "" {
		return obj, fmt.Errorf("object kind could not be discoved")
	}

	return obj, nil
}
