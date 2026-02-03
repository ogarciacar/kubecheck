package sdk_test

import (
	"testing"

	"github.com/ogarciacar/kubecheck/sdk"
	"github.com/stretchr/testify/assert"
)

func TestGetHostFreePort(t *testing.T) {

	port1, err := sdk.GetHostFreePort()

	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Greaterf(t, port1, 0, "Expected valid port number, got %d", port1)
	assert.LessOrEqualf(t, port1, 65535, "Expected valid port number, got %d", port1)

	port2, err := sdk.GetHostFreePort()
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.NotEqualf(t, port1, port2, "expected unique free ports, but got identical ports: %d amd %d", port1, port2)
}
