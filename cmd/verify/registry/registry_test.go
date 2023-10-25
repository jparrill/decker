package registry

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
			name:    "Successful PullSecret verification",
			args:    []string{"registry", "--authfile", "testdata/testdata-good.json", "--url", "localhost:5005"},
			wantErr: false,
		},
		{
			name:    "Connectivity error",
			args:    []string{"registry", "--authfile", "testdata/testdata-bad.json", "--url", "localhost:5006"},
			wantErr: true,
		},
		{
			name:    "Error, regsistry not found",
			args:    []string{"registry", "--authfile", "testdata/testdata-bad.json", "--url", "quay.io"},
			wantErr: true,
		},
		{
			name:    "Error, bad credentials",
			args:    []string{"registry", "--authfile", "testdata/testdata-bad.json", "--url", "localhost:5005"},
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
