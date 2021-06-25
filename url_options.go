package jsonapi

import "strings"

// URLOptions ...
type URLOptions struct {
	Prefix string
}

// Path takes a raw path and returns a new one ready to be consumed by the
// library according to the given options.
func (o *URLOptions) Path(p string) string {
	path := strings.Trim(p, "/")
	prefix := strings.Trim(o.Prefix, "/")

	path = strings.TrimLeft(path, prefix)

	return "/" + path
}
