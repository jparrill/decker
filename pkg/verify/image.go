package verify

import (
	"fmt"

	dockerclient "github.com/docker/docker/client"
	"github.com/jparrill/decker/pkg/core/check"
	coreImage "github.com/jparrill/decker/pkg/core/image"
	"github.com/openshift/library-go/pkg/image/reference"
)

func prepareTemporaryImage(dCli *dockerclient.Client, auth, registryURL, filepath string) (*reference.DockerImageReference, error) {

	alpineImage := NewVerifyContainerImage(alpineSampleImage, "", filepath, dCli)

	err := alpineImage.GetImage(true)
	if err != nil {
		return nil, err
	}

	// parse image
	ref, err := reference.Parse(alpineSampleImage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image (%s): %w", alpineSampleImage, err)
	}

	ref.Registry = registryURL

	privateImage := NewVerifyContainerImage(ref.String(), auth, filepath, dCli)

	// tag image locally
	if err := privateImage.RetagImage(alpineSampleImage, ref.String()); err != nil {
		return nil, err
	}

	// make sure the images are there
	if err := alpineImage.EnsureSourceImage(); err != nil {
		return nil, err
	}
	if err := privateImage.EnsureSourceImage(); err != nil {
		return nil, err
	}

	return &ref, nil
}

func (ci *ContainerImage) Verify() error {
	var err error

	ps := NewVerifyPullSecret(ci.FilePath, false, ci.Debug)

	ref, err := reference.Parse(ci.URL)
	check.Checker("Container Image Parsed", err)
	if err != nil {
		return fmt.Errorf("Error parsing container image: %w", err)
	}

	reg, err := NewVerifyRegistry(
		ref.Registry,
		ci.FilePath,
		ci.Debug,
	)
	if err != nil {
		return fmt.Errorf("Error generating registry: %w", err)
	}
	reg.PSData = ps.Data.Auths[ref.Registry]

	reg.FillAuthCredentials()
	err = reg.Encode()
	check.Checker("Encoded Credentials", err)
	if err != nil {
		return fmt.Errorf("Error encoding credentials: %w", err)
	}

	ci.Auth = reg.EAuth

	err = ci.EnsureSourceImage()
	check.Checker("Image Status", err)
	if err != nil {
		return fmt.Errorf("Error validating container image: %w", err)
	}

	return nil
}

func NewVerifyContainerImage(url, auth, filePath string, dCli *dockerclient.Client) *ContainerImage {
	ci, err := coreImage.NewContainerImage(url, auth, filePath, dCli)
	if err != nil {
		panic(err)
	}

	return &ContainerImage{
		ContainerImage: *ci,
	}
}
