package api

// Kind represents the kind of the object
type Kind string

// Version represents the version of the object
type Version string

// Object is an interface that every configuration object
// that can be converted, used & stored needs to implement.
type Object interface {
	// GetObjectKind returns the kind of the object.
	GetObjectKind() Kind
	// GetObjectVersion returns the version of the object.
	GetObjectVersion() Version
}
