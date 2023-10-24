package verify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyRegistryPushAndPull(t *testing.T) {
	testCases := []struct {
		name    string
		rge     *RegistryEntry
		wantErr bool
	}{
		{
			name: "Successful push and pull",
			rge: &RegistryEntry{
				ServerAddress: "localhost:5005",
				Username:      "admin",
				Password:      "admin",
			},
			wantErr: false,
		},
		{
			name: "Error bad credentials",
			rge: &RegistryEntry{
				ServerAddress: "localhost:5005",
				Username:      "sampleUser",
				Password:      "admin1234",
			},
			wantErr: true,
		},
		{
			name: "Error cannot connect to registry",
			rge: &RegistryEntry{
				ServerAddress: "localhost:5000",
				Username:      "sampleUser",
				Password:      "admin1234",
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.rge.VerifyRegistryCredentials()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyRegistryCredentials(t *testing.T) {
	testCases := []struct {
		name    string
		rge     *RegistryEntry
		wantErr bool
	}{
		{
			name: "Successful registry login",
			rge: &RegistryEntry{
				ServerAddress: "localhost:5005",
				Username:      "admin",
				Password:      "admin",
			},
			wantErr: false,
		},
		{
			name: "Login error, bad credentials",
			rge: &RegistryEntry{
				ServerAddress: "localhost:5005",
				Username:      "sampleUser",
				Password:      "admin1234",
			},
			wantErr: true,
		},
		{
			name: "Connection error, cannot connect to registry",
			rge: &RegistryEntry{
				ServerAddress: "localhost:5000",
				Username:      "sampleUser",
				Password:      "admin1234",
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.rge.VerifyRegistryCredentials()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	testCases := []struct {
		name    string
		rg      *Registry
		wantErr bool
	}{
		{
			name: "Successful registry verification",
			rg: &Registry{
				URL:      "localhost:5005",
				FilePath: "assets/test.json",
			},
			wantErr: false,
		},
		{
			name: "Error on registry verification",
			rg: &Registry{
				URL:      "localhost:5006",
				FilePath: "assets/test.json",
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.rg.Verify()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
