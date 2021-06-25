package examples_test

import (
	"fmt"

	"github.com/mfcochauxlaberge/jsonapi"
)

func ExampleRange() {
	// The Range function allows pagination to be performed
	// on any type that implements the Collection interface.

	// Here is a simple collection with some items.
	col := jsonapi.Resources{}
	col.Add(jsonapi.Wrap(User{
		ID:       "u1",
		Username: "rob",
	}))
	col.Add(jsonapi.Wrap(User{
		ID:       "u2",
		Username: "arthur",
	}))
	col.Add(jsonapi.Wrap(User{
		ID:       "u2",
		Username: "henry",
	}))
	col.Add(jsonapi.Wrap(User{
		ID:       "u2",
		Username: "thisguy",
	}))

	// The proper parameters are given to Range to get the
	// desired page.
	page := jsonapi.Range(
		&col,                       // Collection
		[]string{"u1", "u2", "u3"}, // Only include resources with thos IDs
		nil,                        // Filter (empty here, see filter example)
		[]string{"username", "id"}, // Order in which fields are considered for sorting
		2,                          // Number of elements in the page
		1,                          // Page number, 1 being the first page
	)

	// To sort resources, the given fields are used in the
	// order they appear. The first field is used first, and
	// the follwing ones are used only when needed to break
	// an equality with the previous field. Missing fields
	// are appended to the slice in alphabetical order, except
	// for "id" which will always be appended last unless it
	// appears in the original slice.

	fmt.Printf("Len: %d\n", page.Len())
	// Output: 11
	// fmt.Printf("%s\n%s\n", page.At(0).GetID(), page.At(1).GetID())
}
