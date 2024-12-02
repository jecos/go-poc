package types

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFieldGetAliasWithNonEmptyAlias(t *testing.T) {
	t.Parallel()
	f := Field{
		Name:  "name",
		Alias: "alias",
	}
	assert.Equal(t, f.GetAlias(), "alias")
}

func TestFieldGetAliasWithEmptyAlias(t *testing.T) {
	t.Parallel()
	f := Field{
		Name: "name",
	}
	assert.Equal(t, f.GetAlias(), "name")
}
