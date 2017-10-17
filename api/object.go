package api

// Kind represents the kind of the object
type Kind string

// Object is an interface that every configuration object
// that can be converted, used & stored needs to implement.
type Object interface {
	// GetObjectKind returns the kind of the object.
	GetObjectKind() Kind
}
