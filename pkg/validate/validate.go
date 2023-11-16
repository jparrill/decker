package validate

import (
	"github.com/docker/docker/client"
	coreImage "github.com/jparrill/decker/pkg/core/image"
)

const (
	ReleaseImageStreamFile   = "release-manifests/image-references"
	ReleaseImageMetadataFile = "release-manifests/0000_50_installer_coreos-bootimages.yaml"
)

type OCPImage struct {
	coreImage.ContainerImage
}

type OCPVersion struct {
	Name        string `json:"name"`
	PullSpec    string `json:"pullSpec"`
	DownloadURL string `json:"downloadURL"`
}

type RegistryClientProvider struct{}

func NewValidateOCPImage(url, auth, filePath string, dCLi *client.Client) *OCPImage {
	ci, err := coreImage.NewContainerImage(url, auth, filePath, dCLi)
	if err != nil {
		panic(err)
	}

	return &OCPImage{
		ContainerImage: *ci,
	}
}
