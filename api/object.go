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

// ListOptions are the options required to list & watch objects.
type ListOptions struct {
	TypeMeta      `json:",inline"`
	LabelSelector map[string]string `json:"labelSelector,omitempty"`
}

// GetObjectKind satisfies Object interface.
func (l ListOptions) GetObjectKind() Kind {
	return l.Kind
}

// GetObjectVersion satisfies Object interface.
func (l ListOptions) GetObjectVersion() Version {
	return l.Version
}

// GetObjectMetadata isn't needed in this kind of object.
func (l ListOptions) GetObjectMetadata() ObjectMeta {
	return NoObjectMeta
}

// DeepCopy satisfies object interface.
func (l ListOptions) DeepCopy() Object {
	copy := l
	return &copy
}

// ObjectMeta is the metadata all the objects should have.
type ObjectMeta struct {
	// ID is the id of the object.
	ID string `json:"id,omitempty"`
	// Labels are key/value pairs related with the object used to identify the object.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations are free key/value pairs related with the object that aren't queryable.
	Annotations map[string]string `json:"annotations,omitempty"`
}

// NoObjectMeta is a shortcut to specify the object is not an object.
var NoObjectMeta = ObjectMeta{}

// Object is an interface that every configuration object
// that can be converted, used & stored needs to implement.
type Object interface {
	// GetObjectKind returns the kind of the object.
	GetObjectKind() Kind
	// GetObjectVersion returns the version of the object.
	GetObjectVersion() Version
	// GetObjectMetadata returns the metadata of the object.
	GetObjectMetadata() ObjectMeta
	// DeepCopy makes a copy of the object.
	DeepCopy() Object
}

// ListMeta is the metadata all the objects lists should have.
type ListMeta struct {
	// Continue if not empty means that there are more objects remaining in the list
	Continue string `json:"continue,omitempty"`
}

// ObjectList is an object that is also a list of objects.
type ObjectList interface {
	Object
	GetListMetadata() ListMeta
	GetItems() []Object
}
