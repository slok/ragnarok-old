package api

import (
	"fmt"
)

// Kind represents the kind of the object
type Kind string

// Version represents the version of the object
type Version string

// TypeMeta is the meta type all the objects should have
type TypeMeta struct {
	// Kind represents the kind of the object
	Kind Kind `json:"kind,omitempty"`
	// Version represents the version of the object
	Version Version `json:"version,omitempty"`
}

func (t TypeMeta) String() string {
	return fmt.Sprintf("%s/%s", t.Version, t.Kind)
}

// GetObjectKind satisfies Object interface.
func (t TypeMeta) GetObjectKind() Kind {
	return t.Kind
}

// GetObjectVersion satisfies Object interface.
func (t TypeMeta) GetObjectVersion() Version {
	return t.Version
}

// Object is an interface that every configuration object
// that can be converted, used & stored needs to implement.
type Object interface {
	// GetObjectKind returns the kind of the object.
	GetObjectKind() Kind
	// GetObjectVersion returns the version of the object.
	GetObjectVersion() Version
}
