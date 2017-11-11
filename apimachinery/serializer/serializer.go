package serializer

import (
	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/log"
)

// Global helpers.
var (
	// JSONSerializerDefault is the default base serializer for json.
	JSONSerializerDefault = NewJSONSerializer(ObjTyper, ObjFactory, log.Base())
	// YAMLSerializerDefault is the default base serializer for yaml.
	YAMLSerializerDefault = NewYAMLSerializer(ObjTyper, ObjFactory, log.Base())
	// PBSerializerDefault is the default base serializer for pb.
	PBSerializerDefault = NewPBSerializer(log.Base())
	// DefaultSerializer is the default serializer for the applicaton.
	DefaultSerializer = JSONSerializerDefault
)

// Serializer has ability to serialize objects back and forward in different byte formats.
type Serializer interface {
	// Encode will take an object and will write the serialized data on the received writer.
	Encode(obj api.Object, out interface{}) error
	// Decode will take raw data in format and return a runtime object.
	Decode(data interface{}) (api.Object, error)
}
