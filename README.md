<div align="center" style="text-align: center;">
  <img src="https://raw.githubusercontent.com/mfcochauxlaberge/jsonapi/master/assets/logo.png" height="120">
  <br>
  <a href="https://github.com/mfcochauxlaberge/jsonapi/actions?query=workflow%3ATest+branch%3Amaster">
    <img src="https://github.com/mfcochauxlaberge/jsonapi/workflows/Test/badge.svg?branch=master">
  </a>
  <a href="https://github.com/mfcochauxlaberge/jsonapi/actions?query=workflow%3ALint+branch%3Amaster">
    <img src="https://github.com/mfcochauxlaberge/jsonapi/workflows/Lint/badge.svg?branch=master">
  </a>
  <a href="https://goreportcard.com/report/github.com/mfcochauxlaberge/jsonapi">
    <img src="https://goreportcard.com/badge/github.com/mfcochauxlaberge/jsonapi">
  </a>
  <a href="https://codecov.io/gh/mfcochauxlaberge/jsonapi">
    <img src="https://img.shields.io/codecov/c/github/mfcochauxlaberge/jsonapi">
  </a>
  <br>
  <a href="https://github.com/mfcochauxlaberge/jsonapi/blob/master/go.mod">
    <img src="https://img.shields.io/badge/go%20version-1.13%2B-%2300acd7">
  </a>
  <a href="https://github.com/mfcochauxlaberge/jsonapi/blob/master/go.mod">
    <img src="https://img.shields.io/github/v/release/mfcochauxlaberge/jsonapi?include_prereleases&sort=semver">
  </a>
  <a href="https://github.com/mfcochauxlaberge/jsonapi/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/mfcochauxlaberge/jsonapi?color=a33">
  </a>
  <a href="https://pkg.go.dev/github.com/mfcochauxlaberge/jsonapi?tab=doc">
    <img src="https://img.shields.io/static/v1?label=doc&message=pkg.go.dev&color=007d9c">
  </a>
</div>

# jsonapi

jsonapi offers a set of tools to build JSON:API compliant services.

The official JSON:API specification can be found at [jsonapi.org/format](http://jsonapi.org/format).

## Features

jsonapi offers the following features:

* Marshaling and unmarshaling of JSON:API URLs and documents
* Structs for handling URLs, documents, resources, collections...
* Schema management
  * It can ensure relationships between types make sense.
  * Very useful for validation when marshaling and unmarshaling.
* Utilities for pagination, sorting, and filtering
  * jsonapi is opiniated when it comes to those features. If you prefer you own strategy fo pagination, sorting, and filtering, it will have to be done manually.
* In-memory data store (`SoftCollection`)
  * It can store resources (anything that implements `Resource`).
  * It can sort, filter, retrieve pages, etc.
  * Enough to build a demo API or use in test suites.
  * Not made for production use.
* Other useful helpers

## State

The library is in **beta** and its API is subject to change until v1 is released.

In terms of features, jsonapi is complete. The work left is polishing and testing the design of current API.

### Roadmap to v1

While anything can happen before a v1 release, the API is stable and no big changes are expected at this moment.

A few tasks are required before committing to the current API:

* Rethink how errors are handled
  * Use the new tools introduced in Go 1.13.
* Simplify the API
  * Remove anything that is redundant or not useful.
* Gather feedback from users
  * The library should be used more on real projects to see of the API is convenient.

## Requirements

The supported versions of Go are the latest patch releases of every minor release starting with Go 1.13.

## Examples

The best way to learn and appreciate this package is to look at the simple examples provided in the `examples/` directory.

## Quick start

The simplest way to start using jsonapi is to use the MarshalDocument and UnmarshalDocument functions.

```go
func MarshalDocument(doc *Document, url *URL) ([]byte, error)
func UnmarshalDocument(payload []byte, schema *Schema) (*Document, error)
```

A struct has to follow certain rules in order to be understood by the library, but interfaces are also provided which let the library avoid the reflect package and be more efficient.

See the following section for more information about how to define structs for this library.

## Concepts

Here are some of the main concepts covered by the library.

### Request

A `Request` represents an HTTP request structured in a format easily readable from a JSON:API point of view.

If you are familiar with the specification, reading the `Request` struct and its fields (`URL`, `Document`, etc) should be straightforward.

### Schema

A `Schema` contains all the schema information for an API, like types, fields, relationships between types, and so on. See `schema.go` and `type.go` for more details.

This is really useful for many uses cases:

* Making sure the schema is coherent
* Validating resources
* Parsing documents and URLs
* And probably many more...

For example, when a request comes in, a `Document` and a `URL` can be created by parsing the request. By providing a schema, the parsing can fail if it finds some errors like a type that does not exist, a field of the wrong kind, etc. After that step, valid data can be assumed.

### Type

A JSON:API type is generally defined with a struct.

There needs to be an ID field of type string. The `api` tag represents the name of the type.

```go
type User struct {
  ID string `json:"id" api:"users"` // ID is mandatory and the api tag sets the type

  // Attributes
  Name string `json:"name" api:"attr"` // attr means it is an attribute
  BornAt time.Time `json:"born-at" api:"attr"`

  // Relationships
  Articles []string `json:"articles" api:"rel,articles"`
}
```

Other fields with the `api` tag (`attr` or `rel`) can be added as attributes or relationships.

#### Attribute

Attributes can be of the following types:

```go
string
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64
bool
time.Time
[]byte
*string
*int, *int8, *int16, *int32, *int64
*uint, *uint8, *uint16, *uint32, *uint64
*bool
*time.Time
*[]byte
```

Using a pointer allows the field to be nil.

#### Relationship

Relationships can be a bit tricky. To-one relationships are defined with a string and to-many relationships are defined with a slice of strings. They contain the IDs of the related resources. The api tag has to take the form of "rel,xxx[,yyy]" where yyy is optional. xxx is the type of the relationship and yyy is the name of the inverse relationship when dealing with a two-way relationship. In the following example, our Article struct defines a relationship named author of type users:

```go
Author string `json:"author" api:"rel,users,articles"`
```

### Wrapper

A struct can be wrapped using the `Wrap` function which returns a pointer to a `Wrapper`. A `Wrapper` implements the `Resource` interface and can be used with this library. Modifying a Wrapper will modify the underlying struct. The resource's type is defined from reflecting on the struct.

```go
user := User{}
wrap := Wrap(&user)
wrap.Set("name", "Mike")
fmt.Printf(wrap.Get("name")) // Output: Mike
fmt.Printf(user.Name) // Output: Mike
```

### SoftResource

A SoftResource is a struct whose type (name, attributes, and relationships) can be modified indefinitely just like its values. When an attribute or a relationship is added, the new value is the zero value of the field type. For example, if you add an attribute named `my-attribute` of type string, then `softresource.Get("my-attribute")` will return an empty string.

```go
sr := SoftResource{}
sr.AddAttr(Attr{
  Name:     "attr",
  Type:     AttrTypeInt,
  Nullable: false,
})
fmt.Println(sr.Get("attr")) // Output: 0
```

Take a look at the `SoftCollection` struct for a similar concept applied to an entire collection of resources.

### URLs

From a raw string that represents a URL, it is possible that create a `SimpleURL` which contains the information stored in the URL in a structure that is easier to handle.

It is also possible to build a `URL` from a `Schema` and a `SimpleURL` which contains additional information taken from the schema. `NewURL` returns an error if the URL does not respect the schema.

## Documentation

Check out the [documentation](https://pkg.go.dev/github.com/mfcochauxlaberge/jsonapi?tab=doc).

The best way to learn how to use it is to look at documentation, the examples, and the code itself.

## Contributing

TODO.