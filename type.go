package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Attribute types
const (
	AttrTypeInvalid = ""
	AttrTypeString  = "string"
	AttrTypeInt     = "int"
	AttrTypeInt8    = "int8"
	AttrTypeInt16   = "int16"
	AttrTypeInt32   = "int32"
	AttrTypeInt64   = "int64"
	AttrTypeUint    = "uint"
	AttrTypeUint8   = "uint8"
	AttrTypeUint16  = "uint16"
	AttrTypeUint32  = "uint32"
	AttrTypeUint64  = "uint64"
	AttrTypeBool    = "bool"
	AttrTypeTime    = "time"
	AttrTypeBytes   = "bytes"
)

// A Type stores all the necessary information about a type as represented in
// the JSON:API specification.
//
// NewFunc stores a function that returns a new Resource of the type defined by
// the object with all the fields and the ID set to their zero values. Users may
// call the New method which returns the result of NewFunc if it is non-nil,
// otherwise it returns a SoftResource based on the type.
//
// New makes sure NewFunc is not nil and then calls it, but does not use any
// kind of locking in the process. Therefore, it is unsafe to set NewFunc and
// call New concurrently.
type Type struct {
	Name    string
	Attrs   map[string]Attr
	Rels    map[string]Rel
	NewFunc func() Resource
}

// AddAttr adds an attributes to the type.
func (t *Type) AddAttr(attr Attr) error {
	// Validation
	if attr.Name == "" {
		return fmt.Errorf("jsonapi: attribute name is empty")
	}

	if GetAttrTypeString(attr.Type, attr.Nullable) == "" {
		return fmt.Errorf("jsonapi: attribute type is invalid")
	}

	// Make sure the name isn't already used
	for i := range t.Attrs {
		if t.Attrs[i].Name == attr.Name {
			return fmt.Errorf("jsonapi: attribute name %q is already used", attr.Name)
		}
	}

	if t.Attrs == nil {
		t.Attrs = map[string]Attr{}
	}

	t.Attrs[attr.Name] = attr

	return nil
}

// RemoveAttr removes an attribute from the type if it exists.
func (t *Type) RemoveAttr(attr string) {
	for i := range t.Attrs {
		if t.Attrs[i].Name == attr {
			delete(t.Attrs, attr)
		}
	}
}

// AddRel adds a relationship to the type.
func (t *Type) AddRel(rel Rel) error {
	// Validation
	if rel.FromName == "" {
		return fmt.Errorf("jsonapi: relationship name is empty")
	}

	if rel.ToType == "" {
		return fmt.Errorf("jsonapi: relationship type is empty")
	}

	// Make sure the name isn't already used
	for i := range t.Rels {
		if t.Rels[i].FromName == rel.FromName {
			return fmt.Errorf("jsonapi: relationship name %q is already used", rel.FromName)
		}
	}

	if t.Rels == nil {
		t.Rels = map[string]Rel{}
	}

	t.Rels[rel.FromName] = rel

	return nil
}

// RemoveRel removes a relationship from the type if it exists.
func (t *Type) RemoveRel(rel string) {
	for i := range t.Rels {
		if t.Rels[i].FromName == rel {
			delete(t.Rels, rel)
		}
	}
}

// Fields returns a list of the names of all the fields (attributes and
// relationships) in the type.
func (t *Type) Fields() []string {
	fields := make([]string, 0, len(t.Attrs)+len(t.Rels))
	for i := range t.Attrs {
		fields = append(fields, t.Attrs[i].Name)
	}

	for i := range t.Rels {
		fields = append(fields, t.Rels[i].FromName)
	}

	sort.Strings(fields)

	return fields
}

// New calls the NewFunc field and returns the result Resource object.
//
// If NewFunc is nil, it returns a *SoftResource with its Type field set to the
// value of the receiver.
func (t *Type) New() Resource {
	if t.NewFunc != nil {
		return t.NewFunc()
	}

	return &SoftResource{Type: t}
}

// Equal returns true if both types have the same name, attributes,
// relationships. NewFunc is ignored.
func (t Type) Equal(typ Type) bool {
	t.NewFunc = nil
	typ.NewFunc = nil

	return reflect.DeepEqual(t, typ)
}

// Copy deeply copies the receiver and returns the result.
func (t Type) Copy() Type {
	ctyp := Type{
		Name:  t.Name,
		Attrs: map[string]Attr{},
		Rels:  map[string]Rel{},
	}

	for name, attr := range t.Attrs {
		ctyp.Attrs[name] = attr
	}

	for name, rel := range t.Rels {
		ctyp.Rels[name] = rel
	}

	ctyp.NewFunc = t.NewFunc

	return ctyp
}

// Attr represents a resource attribute.
type Attr struct {
	Name     string
	Type     string
	Nullable bool
}

// UnmarshalToType unmarshals the data into a value of the type represented by
// the attribute and returns it.
func (a Attr) UnmarshalToType(data []byte) (interface{}, error) {
	if a.Nullable && string(data) == "nil" {
		return nil, nil
	}

	var (
		v   interface{}
		err error
	)

	switch a.Type {
	case AttrTypeString:
		var s string
		err = json.Unmarshal(data, &s)

		if a.Nullable {
			v = &s
		} else {
			v = s
		}
	case AttrTypeInt:
		v, err = strconv.Atoi(string(data))

		if a.Nullable {
			n := v.(int)
			v = &n
		} else {
			v = v.(int)
		}
	case AttrTypeInt8:
		v, err = strconv.Atoi(string(data))

		if a.Nullable {
			n := int8(v.(int))
			v = &n
		} else {
			v = int8(v.(int))
		}
	case AttrTypeInt16:
		v, err = strconv.Atoi(string(data))

		if a.Nullable {
			n := int16(v.(int))
			v = &n
		} else {
			v = int16(v.(int))
		}
	case AttrTypeInt32:
		v, err = strconv.Atoi(string(data))

		if a.Nullable {
			n := int32(v.(int))
			v = &n
		} else {
			v = int32(v.(int))
		}
	case AttrTypeInt64:
		v, err = strconv.Atoi(string(data))

		if a.Nullable {
			n := int64(v.(int))
			v = &n
		} else {
			v = int64(v.(int))
		}
	case AttrTypeUint:
		v, err = strconv.ParseUint(string(data), 10, 64)

		if a.Nullable {
			n := uint(v.(uint64))
			v = &n
		} else {
			v = uint(v.(uint64))
		}
	case AttrTypeUint8:
		v, err = strconv.ParseUint(string(data), 10, 8)

		if a.Nullable {
			n := uint8(v.(uint64))
			v = &n
		} else {
			v = uint8(v.(uint64))
		}
	case AttrTypeUint16:
		v, err = strconv.ParseUint(string(data), 10, 16)

		if a.Nullable {
			n := uint16(v.(uint64))
			v = &n
		} else {
			v = uint16(v.(uint64))
		}
	case AttrTypeUint32:
		v, err = strconv.ParseUint(string(data), 10, 32)

		if a.Nullable {
			n := uint32(v.(uint64))
			v = &n
		} else {
			v = uint32(v.(uint64))
		}
	case AttrTypeUint64:
		v, err = strconv.ParseUint(string(data), 10, 64)

		if a.Nullable {
			n := v.(uint64)
			v = &n
		} else {
			v = v.(uint64)
		}
	case AttrTypeBool:
		var b bool
		if string(data) == "true" {
			b = true
		} else if string(data) != "false" {
			err = errors.New("boolean is not true or false")
		}

		v = b

		if a.Nullable {
			v = &b
		}
	case AttrTypeTime:
		var t time.Time
		err = json.Unmarshal(data, &t)
		v = t

		if a.Nullable {
			v = &t
		}
	case AttrTypeBytes:
		s := make([]byte, len(data))
		err := json.Unmarshal(data, &s)

		if err != nil {
			panic(err)
		}

		if a.Nullable {
			v = &s
		} else {
			v = s
		}
	default:
		err = errors.New("attribute is of invalid or unknown type")
	}

	if err != nil {
		return nil, NewErrInvalidFieldValueInBody(
			a.Name,
			string(data),
			GetAttrTypeString(a.Type, a.Nullable),
		)
	}

	return v, nil
}

// Rel represents a resource relationship.
type Rel struct {
	FromType string
	FromName string
	ToOne    bool
	ToType   string
	ToName   string
	FromOne  bool
}

// Invert returns the inverse relationship of r.
func (r *Rel) Invert() Rel {
	return Rel{
		FromType: r.ToType,
		FromName: r.ToName,
		ToOne:    r.FromOne,
		ToType:   r.FromType,
		ToName:   r.FromName,
		FromOne:  r.ToOne,
	}
}

// Normalize inverts the relationship if necessary in order to have it in the
// right direction and returns the result.
//
// This is the form stored in Schema.Rels.
func (r *Rel) Normalize() Rel {
	from := r.FromType + r.FromName
	to := r.ToType + r.ToName

	if from < to || r.ToName == "" {
		return *r
	}

	return r.Invert()
}

// String builds and returns the name of the receiving Rel.
//
// r.Normalize is always called.
func (r Rel) String() string {
	r = r.Normalize()

	id := r.FromType + "_" + r.FromName
	if r.ToName != "" {
		id += "_" + r.ToType + "_" + r.ToName
	}

	return id
}

// GetAttrType returns the attribute type as an int (see constants) and a
// boolean that indicates whether the attribute can be null or not.
func GetAttrType(t string) (int, bool) {
	nullable := strings.HasPrefix(t, "*")

	if nullable {
		t = t[1:]
	}

	switch t {
	case "string":
		return AttrTypeString, nullable
	case "int":
		return AttrTypeInt, nullable
	case "int8":
		return AttrTypeInt8, nullable
	case "int16":
		return AttrTypeInt16, nullable
	case "int32":
		return AttrTypeInt32, nullable
	case "int64":
		return AttrTypeInt64, nullable
	case "uint":
		return AttrTypeUint, nullable
	case "uint8":
		return AttrTypeUint8, nullable
	case "uint16":
		return AttrTypeUint16, nullable
	case "uint32":
		return AttrTypeUint32, nullable
	case "uint64":
		return AttrTypeUint64, nullable
	case "bool":
		return AttrTypeBool, nullable
	case "time.Time":
		return AttrTypeTime, nullable
	case "[]uint8":
		return AttrTypeBytes, nullable
	default:
		return AttrTypeInvalid, false
	}
}

// // GetAttrTypeString returns the name of the attribute type specified by t (see
// // constants) and nullable.
// func GetAttrTypeString(t int, nullable bool) string {
// 	str := ""

// 	switch t {
// 	case AttrTypeString:
// 		str = "string"
// 	case AttrTypeInt:
// 		str = "int"
// 	case AttrTypeInt8:
// 		str = "int8"
// 	case AttrTypeInt16:
// 		str = "int16"
// 	case AttrTypeInt32:
// 		str = "int32"
// 	case AttrTypeInt64:
// 		str = "int64"
// 	case AttrTypeUint:
// 		str = "uint"
// 	case AttrTypeUint8:
// 		str = "uint8"
// 	case AttrTypeUint16:
// 		str = "uint16"
// 	case AttrTypeUint32:
// 		str = "uint32"
// 	case AttrTypeUint64:
// 		str = "uint64"
// 	case AttrTypeBool:
// 		str = "bool"
// 	case AttrTypeTime:
// 		str = "time.Time"
// 	case AttrTypeBytes:
// 		str = "[]uint8"
// 	default:
// 		str = ""
// 	}

// 	if nullable {
// 		return "*" + str
// 	}

// 	return str
// }

// GetZeroValue returns the zero value of the attribute type represented by the
// specified int (see constants).
//
// If nullable is true, the returned value is a nil pointer.
func GetZeroValue(t int, nullable bool) interface{} {
	switch t {
	case AttrTypeString:
		if nullable {
			return (*string)(nil)
		}

		return ""
	case AttrTypeInt:
		if nullable {
			return (*int)(nil)
		}

		return int(0)
	case AttrTypeInt8:
		if nullable {
			return (*int8)(nil)
		}

		return int8(0)
	case AttrTypeInt16:
		if nullable {
			return (*int16)(nil)
		}

		return int16(0)
	case AttrTypeInt32:
		if nullable {
			return (*int32)(nil)
		}

		return int32(0)
	case AttrTypeInt64:
		if nullable {
			return (*int64)(nil)
		}

		return int64(0)
	case AttrTypeUint:
		if nullable {
			return (*uint)(nil)
		}

		return uint(0)
	case AttrTypeUint8:
		if nullable {
			return (*uint8)(nil)
		}

		return uint8(0)
	case AttrTypeUint16:
		if nullable {
			return (*uint16)(nil)
		}

		return uint16(0)
	case AttrTypeUint32:
		if nullable {
			return (*uint32)(nil)
		}

		return uint32(0)
	case AttrTypeUint64:
		if nullable {
			return (*uint64)(nil)
		}

		return uint64(0)
	case AttrTypeBool:
		if nullable {
			return (*bool)(nil)
		}

		return false
	case AttrTypeTime:
		if nullable {
			return (*time.Time)(nil)
		}

		return time.Time{}
	case AttrTypeBytes:
		if nullable {
			return (*[]byte)(nil)
		}

		return []byte{}
	default:
		return nil
	}
}
