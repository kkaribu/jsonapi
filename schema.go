package jsonapi

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Attribute types
const (
	AttrTypeInvalid = iota
	AttrTypeString
	AttrTypeInt
	AttrTypeInt8
	AttrTypeInt16
	AttrTypeInt32
	AttrTypeInt64
	AttrTypeUint
	AttrTypeUint8
	AttrTypeUint16
	AttrTypeUint32
	AttrTypeUint64
	AttrTypeBool
	AttrTypeTime
)

// A Schema contains a list of types. It makes sure that each type is
// valid and unique.
//
// Check can be used to validate the relationships between the types.
type Schema struct {
	Types []Type
}

// AddType adds a type to the schema.
func (s *Schema) AddType(typ Type) error {
	// Validation
	if typ.Name == "" {
		return errors.New("jsonapi: type name is empty")
	}

	// Make sure the name isn't already used
	for i := range s.Types {
		if s.Types[i].Name == typ.Name {
			return fmt.Errorf("jsonapi: type name %s is already used", typ.Name)
		}
	}

	s.Types = append(s.Types, typ)

	return nil
}

// RemoveType removes a type from the schema.
func (s *Schema) RemoveType(typ string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			s.Types = append(s.Types[0:i], s.Types[i+1:]...)
		}
	}

	return nil
}

// AddAttr adds an attribute to the specified type.
func (s *Schema) AddAttr(typ string, attr Attr) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].AddAttr(attr)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// RemoveAttr removes an attribute from the specified type.
func (s *Schema) RemoveAttr(typ string, attr string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].RemoveAttr(attr)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// AddRel adds a relationship to the specified type.
func (s *Schema) AddRel(typ string, rel Rel) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].AddRel(rel)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// RemoveRel removes a relationship from the specified type.
func (s *Schema) RemoveRel(typ string, rel string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].RemoveRel(rel)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// HasType returns a boolean indicating whether a type has the specified name
// or not.
func (s *Schema) HasType(name string) bool {
	for i := range s.Types {
		if s.Types[i].Name == name {
			return true
		}
	}
	return false
}

// GetType returns the type associated with the speficied name.
//
// A boolean indicates whether a type was found or not.
func (s *Schema) GetType(name string) (Type, bool) {
	for _, typ := range s.Types {
		if typ.Name == name {
			return typ, true
		}
	}
	return Type{}, false
}

// GetResource returns a resource of type SoftResource with the specified
// type. All fields are set to their zero values.
func (s *Schema) GetResource(name string) Resource {
	typ, ok := s.GetType(name)
	if ok {
		return NewSoftResource(typ, nil)
	}
	return nil
}

// Check checks the integrity of all the relationships between the types
// and returns all the errors that were found.
func (s *Schema) Check() []error {
	var (
		ok   bool
		errs = []error{}
	)

	// Check the inverse relationships
	for _, typ := range s.Types {
		// Relationships
		for _, rel := range typ.Rels {
			var targetType Type

			// Does the relationship point to a type that exists?
			if targetType, ok = s.GetType(rel.Type); !ok {
				errs = append(errs, fmt.Errorf(
					"jsonapi: the target type of relationship %s of type %s does not exist",
					rel.Name,
					typ.Name,
				))
			}

			// Inverse relationship (if relevant)
			if rel.InverseName != "" {
				// Is the inverse relationship type the same as its type name?
				if rel.InverseType != typ.Name {
					errs = append(errs, fmt.Errorf(
						"jsonapi: the inverse type of relationship %s should its type's name (%s, not %s)",
						rel.Name,
						typ.Name,
						rel.InverseType,
					))
				}

				// Do both relationships (current and inverse) point to each other?
				var found bool
				for _, invRel := range targetType.Rels {
					if rel.Name == invRel.InverseName && rel.InverseName == invRel.Name {
						found = true
					}
				}
				if !found {
					errs = append(errs, fmt.Errorf(
						"jsonapi: relationship %s of type %s and its inverse do not point each other",
						rel.Name,
						typ.Name,
					))
				}
			}

		}
	}

	return errs
}

// GetAttrType returns the attribute type as an int (see constants) and
// a boolean that indicates whether the attribute can be null or not.
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
	default:
		return AttrTypeInvalid, false
	}
}

// GetAttrTypeString return the name of the attribute type specified
// by an int (see constants) and a boolean that indicates whether the
// value can be null or not.
func GetAttrTypeString(t int, nullable bool) string {
	str := ""
	switch t {
	case AttrTypeString:
		str = "string"
	case AttrTypeInt:
		str = "int"
	case AttrTypeInt8:
		str = "int8"
	case AttrTypeInt16:
		str = "int16"
	case AttrTypeInt32:
		str = "int32"
	case AttrTypeInt64:
		str = "int64"
	case AttrTypeUint:
		str = "uint"
	case AttrTypeUint8:
		str = "uint8"
	case AttrTypeUint16:
		str = "uint16"
	case AttrTypeUint32:
		str = "uint32"
	case AttrTypeUint64:
		str = "uint64"
	case AttrTypeBool:
		str = "bool"
	case AttrTypeTime:
		str = "time.Time"
	default:
		str = ""
	}
	if nullable {
		return "*" + str
	}
	return str
}

// GetZeroValue returns the zero value of the attribute type represented
// by the specified int (see constants).
//
// If null is true, the returned value is a nil pointer.
func GetZeroValue(t int, null bool) interface{} {
	switch t {
	case AttrTypeString:
		if null {
			var np *string
			return np
		}
		return ""
	case AttrTypeInt:
		if null {
			var np *int
			return np
		}
		return int(0)
	case AttrTypeInt8:
		if null {
			var np *int8
			return np
		}
		return int8(0)
	case AttrTypeInt16:
		if null {
			var np *int16
			return np
		}
		return int16(0)
	case AttrTypeInt32:
		if null {
			var np *int32
			return np
		}
		return int32(0)
	case AttrTypeInt64:
		if null {
			var np *int64
			return np
		}
		return int64(0)
	case AttrTypeUint:
		if null {
			var np *uint
			return np
		}
		return uint(0)
	case AttrTypeUint8:
		if null {
			var np *uint8
			return np
		}
		return uint8(0)
	case AttrTypeUint16:
		if null {
			var np *uint16
			return np
		}
		return uint16(0)
	case AttrTypeUint32:
		if null {
			var np *uint32
			return np
		}
		return uint32(0)
	case AttrTypeUint64:
		if null {
			var np *uint64
			return np
		}
		return uint64(0)
	case AttrTypeBool:
		if null {
			var np *bool
			return np
		}
		return false
	case AttrTypeTime:
		if null {
			var np *time.Time
			return np
		}
		return time.Time{}
	default:
		return nil
	}
}
