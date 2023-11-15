package verify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	staticPsFile = "assets/test.json"
)

func TestVerifyRegistryPushAndPull(t *testing.T) {
	testCases := []struct {
		name          string
		serverAddress string
		username      string
		password      string
		wantErr       bool
	}{
		{
			name:          "Successful push and pull",
			serverAddress: "localhost:5005",
			username:      "admin",
			password:      "admin",
			wantErr:       false,
		},
		{
			name:          "Error bad credentials",
			serverAddress: "localhost:5005",
			username:      "sampleUser",
			password:      "admin1234",
			wantErr:       true,
		},
		{
			name:          "Error cannot connect to registry",
			serverAddress: "localhost:5000",
			username:      "sampleUser",
			password:      "admin1234",
			wantErr:       true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reg := NewVerifyRegistry(tc.serverAddress, staticPsFile, false)
			reg.PSData.Username = tc.username
			reg.PSData.Password = tc.password
			reg.FilePath = staticPsFile
			err := reg.VerifyRegistryPushAndPull()
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
		name          string
		serverAddress string
		username      string
		password      string
		wantErr       bool
	}{
		{
			name:          "Successful registry login",
			serverAddress: "localhost:5005",
			username:      "admin",
			password:      "admin",
			wantErr:       false,
		},
		{
			name:          "Login error, bad credentials",
			serverAddress: "localhost:5005",
			username:      "sampleUser",
			password:      "admin1234",
			wantErr:       true,
		},
		{
			name:          "Connection error, cannot connect to registry",
			serverAddress: "localhost:5000",
			username:      "sampleUser",
			password:      "admin1234",
			wantErr:       true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reg := NewVerifyRegistry(tc.serverAddress, "", false)
			reg.PSData.Username = tc.username
			reg.PSData.Password = tc.password
			err := reg.VerifyRegistryCredentials()
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
		regURL  string
		regPS   string
		wantErr bool
	}{
		{
			name:    "Successful registry verification",
			regURL:  "localhost:5005",
			regPS:   "assets/test.json",
			wantErr: false,
		},
		{
			name:    "Error on registry verification",
			regURL:  "localhost:5006",
			regPS:   "assets/test.json",
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reg := NewVerifyRegistry(tc.regURL, tc.regPS, false)
			err := reg.Verify()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
