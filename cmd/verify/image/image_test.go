package image

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVerifyCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "Successful container image verification",
			args:    []string{"image", "--authfile", "testdata.json", "--url", "localhost:5005/libpod/alpine:latest"},
			wantErr: false,
		},
		{
			name:    "Error, missing url flag",
			args:    []string{"image", "--authfile", "testdata.json", "--url"},
			wantErr: true,
		},
		{
			name:    "Error, container image not found in registry",
			args:    []string{"image", "--authfile", "testdata.json", "--url", "quay.io/libpod/alpina:latest"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up command flags
			cmd := NewVerifyCommand()
			cmd.SetArgs(tt.args)
			cmd.SetOut(os.Stdout)
			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
