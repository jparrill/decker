// BEGIN: xz3y6a8b9cde
package verify

import (
	"context"
	"errors"
	"testing"

	dockerregistrytype "github.com/docker/docker/api/types/registry"
	"github.com/stretchr/testify/assert"
)

type mockDockerClient struct {
	registryLoginFunc func(ctx context.Context, auth dockerregistrytype.AuthConfig) (dockerregistrytype.AuthenticateOKBody, error)
}

func (m *mockDockerClient) RegistryLogin(ctx context.Context, auth dockerregistrytype.AuthConfig) (dockerregistrytype.AuthenticateOKBody, error) {
	return m.registryLoginFunc(ctx, auth)
}

func TestVerifyRegistryCredentials(t *testing.T) {
	tests := []struct {
		name             string
		registryEntry    *RegistryEntry
		mockDockerClient *mockDockerClient
		expectedError    error
	}{
		{
			name: "successful login",
			registryEntry: &RegistryEntry{
				ServerAddress: "https://example.com",
				Username:      "user",
				Password:      "password",
			},
			mockDockerClient: &mockDockerClient{
				registryLoginFunc: func(ctx context.Context, auth dockerregistrytype.AuthConfig) (dockerregistrytype.AuthenticateOKBody, error) {
					return dockerregistrytype.AuthenticateOKBody{}, nil
				},
			},
			expectedError: nil,
		},
		{
			name: "failed login",
			registryEntry: &RegistryEntry{
				ServerAddress: "https://example.com",
				Username:      "user",
				Password:      "password",
			},
			mockDockerClient: &mockDockerClient{
				registryLoginFunc: func(ctx context.Context, auth dockerregistrytype.AuthConfig) (dockerregistrytype.AuthenticateOKBody, error) {
					return dockerregistrytype.AuthenticateOKBody{}, errors.New("failed to login")
				},
			},
			expectedError: errors.New("failed to login"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//dockerClient := tt.mockDockerClient
			err := tt.registryEntry.VerifyRegistryCredentials()
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

// END: xz3y6a8b9cde
