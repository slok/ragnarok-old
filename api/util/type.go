package util

import (
	"strings"

	"github.com/slok/ragnarok/api"
)

const (
	joinChar  = "/"
	apiPrefix = "/api"
)

// GetFullIDFromType will return the full id of a type and an id.
func GetFullIDFromType(t api.TypeMeta, id string) string {
	strs := []string{
		string(t.Version),
		string(t.Kind),
		id,
	}
	return strings.Join(strs, joinChar)
}

// GetFullID will return the full id of an object.
func GetFullID(obj api.Object) string {
	strs := []string{
		string(obj.GetObjectVersion()),
		string(obj.GetObjectKind()),
		obj.GetObjectMetadata().ID,
	}
	return strings.Join(strs, joinChar)
}

// GetFullType will return the object id of an object.
func GetFullType(obj api.Object) string {
	strs := []string{string(obj.GetObjectVersion()), string(obj.GetObjectKind())}
	return strings.Join(strs, joinChar)
}

// GetTypeAPIPath returns the object API path/route.
func GetTypeAPIPath(t api.TypeMeta) string {
	strs := []string{apiPrefix, string(t.GetObjectVersion()), string(t.GetObjectKind())}
	return strings.Join(strs, joinChar)
}

// SplitFullID will split a full id and return a type and an id
func SplitFullID(fullID string) (api.TypeMeta, string) {
	t := api.TypeMeta{}
	spl := strings.Split(fullID, joinChar)
	if len(spl) != 4 {
		return t, ""
	}

	version := strings.Join([]string{spl[0], spl[1]}, joinChar)
	t.Version = api.Version(version)
	t.Kind = api.Kind(spl[2])
	id := spl[3]

	return t, id
}
