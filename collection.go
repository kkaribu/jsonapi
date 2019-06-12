package jsonapi

// A Collection can manage a set of ordered resources of the same type.
type Collection interface {
	// Type returns the name of the resources' type.
	Type() string
	// Len return the number of resources in the collection.
	Len() int
	// Elem returns the resource at index i.
	Elem(i int) Resource
	// Add adds a resource in the collection.
	Add(r Resource)

	// UnmarshalJSON unmarshals the bytes that represent a collection
	// of resources into the struct that implements the interface.
	UnmarshalJSON([]byte) error
}
