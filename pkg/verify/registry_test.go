package verify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyRegistryPushAndPull(t *testing.T) {
	// Create a new RegistryEntry with test credentials
	rge := &RegistryEntry{
		ServerAddress: "test.registry.com",
		Username:      "testuser",
		Password:      "testpassword",
	}

	// Verify the registry credentials
	err := rge.VerifyRegistryCredentials()
	assert.NoError(t, err)

	// Verify the registry push and pull
	err = rge.VerifyRegistryPushAndPull()
	assert.NoError(t, err)
}
