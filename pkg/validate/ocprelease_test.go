package validate

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		SrcImage string
		DstImage string
		wantErr  bool
	}{
		{
			name:     "Given valid images, should return no error",
			SrcImage: "quay.io/openshift-release-dev/ocp-release:4.14.1-x86_64",
			DstImage: "quay.io/openshift-release-dev/ocp-release:4.14.1-x86_64",
			wantErr:  false,
		},
		{
			name:     "invalid source image",
			SrcImage: "quay.io/openshift-release-dev/ocp-release:4.14.1",
			DstImage: "quay.io/openshift-release-dev/ocp-release:4.14.1-x86_64",
			wantErr:  true,
		},
		{
			name:     "invalid destination image",
			SrcImage: "quay.io/openshift-release-dev/ocp-release:4.14.1-x86_64",
			DstImage: "quay.io/openshift-release-dev/ocp-release:4.14.1",
			wantErr:  true,
		},
		{
			name:     "different images",
			SrcImage: "quay.io/openshift-release-dev/ocp-release:4.14.1-x86_64",
			DstImage: "quay.io/openshift-release-dev/ocp-release:4.14.0-x86_64",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcImage := NewValidateOCPImage(tt.SrcImage, "", "", nil)
			dstImage := NewValidateOCPImage(tt.DstImage, "", "", nil)

			if err := Validate(srcImage, dstImage); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
