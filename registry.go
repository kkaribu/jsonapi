package jsonapi

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// Registry ...
type Registry struct {
	sync.Mutex

	Types map[string]Type
}

// NewRegistry ...
func NewRegistry() *Registry {
	return &Registry{
		Types: map[string]Type{},
	}
}

// RegisterType checks and registers the provided value as a type.
func (r *Registry) RegisterType(v interface{}) {
	r.Lock()
	defer r.Unlock()

	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Struct {
		value = reflect.New(value.Type())
	}

	if value.Kind() != reflect.Ptr {
		panic("jsonapi: RegisterType requires a struct or a pointer to a struct")
	}

	if value.Elem().Kind() != reflect.Struct {
		panic("jsonapi: RegisterType requires a struct or a pointer to a struct")
	}

	err := CheckType(value.Elem().Interface())
	if err != nil {
		panic(err)
	}

	value = value.Elem()

	// Get ID field
	idField, _ := value.Type().FieldByName("ID")

	// Get name
	resType := idField.Tag.Get("api")

	fields := []string{}

	// Get attributes
	attrs := map[string]Attr{}
	for i := 0; i < value.NumField(); i++ {
		sf := value.Type().Field(i)

		if sf.Tag.Get("api") == "attr" {
			var n string
			if n = sf.Tag.Get("json"); n == "" {
				n = sf.Name
			}

			def := sql.NullString{}
			def.String = ""

			typ, null := GetAttrType(sf.Type.String())
			attrs[n] = Attr{
				Name: n,
				Type: typ,
				Null: null,
			}

			fields = append(fields, n)
		}
	}

	// Get relationships
	rels := map[string]Rel{}
	for i := 0; i < value.NumField(); i++ {
		sf := value.Type().Field(i)

		if strings.Contains(sf.Tag.Get("api"), "rel,") {
			var t, i string
			if s := strings.Split(sf.Tag.Get("api"), ","); len(s) >= 2 {
				t = s[1]

				if len(s) == 3 {
					i = s[2]
				}
			}

			var toOne bool
			if sf.Type.String() == "string" {
				toOne = true
			} else if sf.Type.String() == "[]string" {
				toOne = false
			}

			var n string
			if n = sf.Tag.Get("json"); n == "" {
				n = sf.Name
			}

			rels[n] = Rel{
				Name:         n,
				Type:         t,
				ToOne:        toOne,
				InverseName:  i,
				InverseType:  resType,
				InverseToOne: false, // should be set in Check()
			}

			fields = append(fields, n)
		}
	}

	if _, ok := r.Types[resType]; ok {
		panic("jsonapi: type with same name already exists")
	}

	r.Types[resType] = Type{
		Name:  resType,
		Attrs: attrs,
		Rels:  rels,
	}
}

// Check ...
func (r *Registry) Check() []error {
	errs := []error{}

	// Check and set the inverse relationships
	for t, typ := range r.Types {
		for re, rel := range typ.Rels {
			if _, ok := r.Types[rel.Type]; !ok {
				errs = append(errs, fmt.Errorf("jsonapi: the target type of relationship %s of type %s does not exist", rel.Name, typ.Name))
			}

			if rel.InverseName != "" {
				if invRel, ok := r.Types[rel.Type].Rels[rel.InverseName]; !ok {
					errs = append(errs, fmt.Errorf("jsonapi: the inverse of relationship %s of type %s does not exist", rel.Name, typ.Name))
				} else {
					rel.InverseToOne = invRel.ToOne
					// rel.InverseToMany = invRel.ToMany
				}

				typ.Rels[re] = rel
			}
		}

		r.Types[t] = typ
	}

	return errs
}

// Resource ...
func (r *Registry) Resource(name string) Resource {
	if t, ok := r.Types[name]; ok {
		return NewSoftResource(t, nil)
	}

	panic(fmt.Sprintf(`jsonapi: type "%s" not found`, name))
}

// Collection ...
func (r *Registry) Collection(name string) Collection {
	if t, ok := r.Types[name]; ok {
		r := NewSoftResource(t, nil)
		col := WrapCollection(r)
		return col
	}

	panic(fmt.Sprintf(`jsonapi: type "%s" not found`, name))
}
