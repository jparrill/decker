package verify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
	"github.com/openshift/library-go/pkg/image/reference"
)

func (ci *ContainerImage) GetImage() error {
	out, err := ci.DClient.ImagePull(context.Background(), ci.ImageURL, types.ImagePullOptions{
		RegistryAuth: ci.Auth,
	})
	if err != nil {
		return fmt.Errorf("Error grabbing container image %s: %v", ci.ImageURL, err)
	}

	defer out.Close()

	if debug {
		_, err := io.Copy(os.Stdout, out)
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", ci.ImageURL, err)
		}
	} else {
		var b bytes.Buffer
		_, err := io.Copy(&b, out)
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", ci.ImageURL, err)
		}
	}

	return nil
}

func (ci *ContainerImage) RetagImage(srcImage, destImage string) error {
	err := ci.DClient.ImageTag(context.Background(), alpineSampleImage, destImage)
	if err != nil {
		return fmt.Errorf("Error re-tagging public container image from %s to %s: %v", alpineSampleImage, destImage, err)
	}

	return nil
}

func (ci *ContainerImage) EnsureSourceImage() error {
	if ci.Debug {
		fmt.Printf("Verifying Image: %s\n", ci.ImageURL)
	}

	filters := filters.NewArgs()
	filters.Add("reference", ci.ImageURL)

	options := types.ImageListOptions{
		All:     true,
		Filters: filters,
	}

	images, err := ci.DClient.ImageList(context.Background(), options)
	if err != nil {
		return fmt.Errorf("error querying the local images: %v", err)
	}

	for _, image := range images {
		if image.RepoTags[0] != ci.ImageURL {
			return fmt.Errorf("The container image %s does not exists localy: %v", ci.ImageURL, err)
		}
	}

	return nil
}

func prepareTemporaryImage(dCli *dockerclient.Client, auth, registryURL string) (*reference.DockerImageReference, error) {

	// Pull source image
	alpineImage := &ContainerImage{
		DClient:  dCli,
		ImageURL: alpineSampleImage,
		Auth:     "",
	}

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

	privateImage := &ContainerImage{
		DClient:  dCli,
		ImageURL: ref.String(),
		Auth:     auth,
	}

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
