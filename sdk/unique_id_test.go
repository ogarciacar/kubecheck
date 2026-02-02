package sdk_test

import (
	"testing"

	"github.com/ogarciacar/kubecheck/sdk"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUniqueID(t *testing.T) {

	// act
	id1 := sdk.GenerateUniqueID()
	id2 := sdk.GenerateUniqueID()

	// assert
	assert.Lenf(t, id1, 8, "Expected length of 8, but got %d", len(id1))
	assert.Lenf(t, id2, 8, "Expected length of 8, but got %d", len(id2))
	assert.NotEqualf(t, id1, id2, "Expected unique IDs, but got identical IDs: %s and %s", id1, id2)
}
