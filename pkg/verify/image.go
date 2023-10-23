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

func getImage(dCli *dockerclient.Client, image, auth string) error {
	out, err := dCli.ImagePull(context.Background(), image, types.ImagePullOptions{
		RegistryAuth: auth,
	})
	if err != nil {
		return fmt.Errorf("Error grabbing container image %s: %v", image, err)
	}

	defer out.Close()

	if debug {
		_, err := io.Copy(os.Stdout, out)
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", image, err)
		}
	} else {
		var b bytes.Buffer
		_, err := io.Copy(&b, out)
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", image, err)
		}
	}

	return nil
}

func retagSourceImage(dCli *dockerclient.Client, srcImage, destImage string) error {
	err := dCli.ImageTag(context.Background(), alpineSampleImage, destImage)
	if err != nil {
		return fmt.Errorf("Error re-tagging public container image from %s to %s: %v", alpineSampleImage, destImage, err)
	}

	return nil
}

func ensureSourceImage(dCli *dockerclient.Client, destImage string) error {
	if debug {
		fmt.Printf("Verifying Image: %s\n", destImage)
	}

	filters := filters.NewArgs()
	filters.Add("reference", destImage)

	options := types.ImageListOptions{
		All:     true,
		Filters: filters,
	}

	images, err := dCli.ImageList(context.Background(), options)
	if err != nil {
		return fmt.Errorf("error querying the local images: %v", err)
	}

	for _, image := range images {
		if image.RepoTags[0] != destImage {
			return fmt.Errorf("The container image %s does not exists localy: %v", destImage, err)
		}
	}

	return nil
}

func prepareTemporaryImage(dCli *dockerclient.Client, auth, registryURL string) (*reference.DockerImageReference, error) {

	// Pull source image
	err := getImage(dCli, alpineSampleImage, auth)
	if err != nil {
		return nil, err
	}

	// parse image
	ref, err := reference.Parse(alpineSampleImage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image (%s): %w", alpineSampleImage, err)
	}

	// tag image locally
	ref.Registry = registryURL
	if err := retagSourceImage(dCli, alpineSampleImage, ref.String()); err != nil {
		return nil, err
	}

	// make sure the images are there
	if err := ensureSourceImage(dCli, alpineSampleImage); err != nil {
		return nil, err
	}
	if err := ensureSourceImage(dCli, ref.String()); err != nil {
		return nil, err
	}

	return &ref, nil
}
