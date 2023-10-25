package pullsecret

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
			args:    []string{"pull-secret", "--authfile", "testdata/testdata-good.json"},
			wantErr: false,
		},
		{
			name:    "Successful PullSecret verification with inspect",
			args:    []string{"pull-secret", "--authfile", "testdata/testdata-good.json", "--inspect"},
			wantErr: false,
		},
		{
			name:    "Error, connectivity error",
			args:    []string{"pull-secret", "--authfile", "testdata/testdata-bad.json", "--inspect"},
			wantErr: true,
		},
		{
			name:    "Multiple errors found",
			args:    []string{"pull-secret", "--authfile", "testdata/testdata-bad-multi.json", "--inspect"},
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
