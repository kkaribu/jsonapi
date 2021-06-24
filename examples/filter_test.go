package examples_test

import (
	"github.com/mfcochauxlaberge/jsonapi"
)

func ExampleFilter() {
	// Here a collection of resources with users in it.
	col := jsonapi.Resources{}
	col.Add(jsonapi.Wrap(&User{
		Username: "bob",
	}))
}
