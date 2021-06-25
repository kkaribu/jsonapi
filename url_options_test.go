package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLOptionsPath(t *testing.T) {
	assert := assert.New(t)

	// Empty
	opts := URLOptions{}

	assert.Equal("/", opts.Path("/"))
	assert.Equal("/type", opts.Path("/type"))
	assert.Equal("/type/id", opts.Path("/type/id"))
	assert.Equal("/type/id/relationship", opts.Path("/type/id/relationship"))
	assert.Equal("/type/id/relationship/id", opts.Path("/type/id/relationship/id"))

	// Slash
	opts = URLOptions{Prefix: "/"}

	assert.Equal("/", opts.Path("/"))
	assert.Equal("/type", opts.Path("/type"))
	assert.Equal("/type/id", opts.Path("/type/id"))
	assert.Equal("/type/id/relationship", opts.Path("/type/id/relationship"))
	assert.Equal("/type/id/relationship/id", opts.Path("/type/id/relationship/id"))

	// Path no slash
	opts = URLOptions{Prefix: "api"}

	assert.Equal("/", opts.Path("/"))
	assert.Equal("/type", opts.Path("/type"))
	assert.Equal("/type/id", opts.Path("/type/id"))
	assert.Equal("/type/id/relationship", opts.Path("/type/id/relationship"))
	assert.Equal("/type/id/relationship/id", opts.Path("/type/id/relationship/id"))
	assert.Equal("/", opts.Path("/api/"))
	assert.Equal("/type", opts.Path("/api/type"))
	assert.Equal("/type/id", opts.Path("/api/type/id"))
	assert.Equal("/type/id/relationship", opts.Path("/api/type/id/relationship"))
	assert.Equal("/type/id/relationship/id", opts.Path("/api/type/id/relationship/id"))
}
