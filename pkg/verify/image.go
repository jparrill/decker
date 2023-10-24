package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
	"github.com/jparrill/decker/pkg/core/check"
	"github.com/openshift/library-go/pkg/image/reference"
)

func (ci *ContainerImage) GetImage() error {
	out, err := ci.DClient.ImagePull(context.Background(), ci.URL, types.ImagePullOptions{
		RegistryAuth: ci.Auth,
	})
	if err != nil {
		return fmt.Errorf("Error grabbing container image %s: %v", ci.URL, err)
	}

	defer out.Close()

	if debug {
		_, err := io.Copy(os.Stdout, out)
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", ci.URL, err)
		}
	} else {
		var b bytes.Buffer
		_, err := io.Copy(&b, out)
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", ci.URL, err)
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

func prepareTemporaryImage(dCli *dockerclient.Client, auth, registryURL string) (*reference.DockerImageReference, error) {

	// Pull source image
	alpineImage := &ContainerImage{
		DClient: dCli,
		URL:     alpineSampleImage,
		Auth:    "",
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
		DClient: dCli,
		URL:     ref.String(),
		Auth:    auth,
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

func (ci *ContainerImage) Verify() error {
	var data AuthsType

	dClient, err := dockerclient.NewClientWithOpts()
	check.Checker("Docker Client Generated", err)
	if err != nil {
		return fmt.Errorf("Error generating docker client: %w", err)
	}
	ci.DClient = dClient

	jsonData, err := os.ReadFile(ci.FilePath)
	check.Checker("Read input file", err)
	if err != nil {
		return fmt.Errorf("Error reading authfile: %w", err)
	}

	err = json.Unmarshal(jsonData, &data)
	check.Checker("Unmarshal JSON file", err)
	if err != nil {
		return fmt.Errorf("Error unmarshaling JSON authfile: %w", err)
	}

	ref, err := reference.Parse(ci.URL)
	check.Checker("Container Image Parsed", err)
	if err != nil {
		return fmt.Errorf("Error parsing container image: %w", err)
	}

	registryEntry := NewRegistryAuth(
		ref.Registry,
		data.Auths[ref.Registry].Username,
		data.Auths[ref.Registry].Password,
		data.Auths[ref.Registry].Auth,
	)

	registryEntry.FillAuthCredentials()
	ci.Auth, err = registryEntry.Encode()
	check.Checker("Encoded Credentials", err)
	if err != nil {
		return fmt.Errorf("Error encoding credentials: %w", err)
	}

	err = ci.EnsureSourceImage()
	check.Checker("Image Status", err)
	if err != nil {
		return fmt.Errorf("Error validating container image: %w", err)
	}

	return nil
}
