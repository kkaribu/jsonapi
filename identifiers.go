package jsonapi

// NewIdentifiers returns an Identifiers object.
//
// t is the type of the identifiers. ids is the set of IDs.
func NewIdentifiers(t string, ids []string) Identifiers {
	identifiers := []Identifier{}

	for _, id := range ids {
		identifiers = append(identifiers, Identifier{
			Type: t,
			ID:   id,
		})
	}

	return identifiers
}

// Identifiers represents a slice of Identifier.
type Identifiers []Identifier

// IDs returns the IDs part of the Identifiers.
func (i Identifiers) IDs() []string {
	ids := []string{}

	for _, id := range i {
		ids = append(ids, id.ID)
	}

	return ids
}

// Identifier represents a resource's type and ID.
type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
