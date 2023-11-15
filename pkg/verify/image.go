package verify

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
	"github.com/jparrill/decker/pkg/core/check"
	coreImage "github.com/jparrill/decker/pkg/core/image"
	"github.com/openshift/library-go/pkg/image/reference"
)

func (ci *ContainerImage) RetagImage(srcImage, destImage string) error {
	err := ci.DClient.ImageTag(context.Background(), alpineSampleImage, destImage)
	if err != nil {
		return fmt.Errorf("Error re-tagging public container image from %s to %s: %v", alpineSampleImage, destImage, err)
	}

	return nil
}

func (ci *ContainerImage) EnsureSourceImage() error {
	if ci.Debug {
		fmt.Printf("Verifying Image: %s\n", ci.URL)
	}

	filters := filters.NewArgs()
	filters.Add("reference", ci.URL)

	options := types.ImageListOptions{
		All:     true,
		Filters: filters,
	}

	images, err := ci.DClient.ImageList(context.Background(), options)
	if err != nil {
		return fmt.Errorf("error querying the local images: %v", err)
	}

	if ci.Debug {
		fmt.Println("Images found:", images)
	}

	if len(images) <= 0 {
		return fmt.Errorf("The container image %s does not exists", ci.URL)
	}

	for _, image := range images {
		if image.RepoTags[0] != ci.URL {
			return fmt.Errorf("The container image %s does not exists", ci.URL)
		}
	}

	return nil
}

func prepareTemporaryImage(dCli *dockerclient.Client, auth, registryURL, filepath string) (*reference.DockerImageReference, error) {

	alpineImage := NewVerifyContainerImage(alpineSampleImage, "", filepath, dCli)

	err := alpineImage.GetImage()
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

	ci.DClient, err = dockerclient.NewClientWithOpts()
	check.Checker("Docker Client Generated", err)
	if err != nil {
		return fmt.Errorf("Error generating docker client: %w", err)
	}

	ref, err := reference.Parse(ci.URL)
	check.Checker("Container Image Parsed", err)
	if err != nil {
		return fmt.Errorf("Error parsing container image: %w", err)
	}

	reg := NewVerifyRegistry(
		ref.Registry,
		ci.FilePath,
		ci.Debug,
	)
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
