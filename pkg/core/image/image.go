package image

import (
	"bufio"
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

type ContainerImage struct {
	DClient   *dockerclient.Client
	URL       string
	FilePath  string
	Auth      string
	TLSVerify bool
	Debug     bool
	Manifest  *types.ImageInspect
	Ref       *reference.DockerImageReference
}

type ContainerImageInterface interface {
	GetImage()
	GetMetadata()
	GetReference()
}

func NewContainerImage(url, auth, filepath string, dCLi *dockerclient.Client) (*ContainerImage, error) {
	ci := &ContainerImage{
		DClient:   dCLi,
		URL:       url,
		FilePath:  filepath,
		Auth:      auth,
		TLSVerify: false,
		Debug:     false,
		Manifest:  nil,
		Ref:       nil,
	}

	err := ci.GetReference()
	if err != nil {
		return nil, fmt.Errorf("Error getting container image reference: %w", err)
	}

	err = ci.SetDockerClient()
	if err != nil {
		return nil, err
	}

	return ci, nil
}

func (ci *ContainerImage) GetReference() error {
	var err error

	ref, err := reference.Parse(ci.URL)
	if err != nil {
		return fmt.Errorf("Error parsing container image: %w", err)
	}

	ci.Ref = &ref

	return nil
}

func (ci *ContainerImage) SetDockerClient() error {
	var err error

	ci.DClient, err = dockerclient.NewClientWithOpts()
	if err != nil {
		return fmt.Errorf("Error generating docker client: %w", err)
	}

	return nil
}

func (ci *ContainerImage) GetMetadata() error {
	var err error
	var manifest types.ImageInspect

	manifest, _, err = ci.DClient.ImageInspectWithRaw(context.Background(), ci.URL)
	if err != nil {
		return fmt.Errorf("Error inspecting docker image %s: %w", ci.URL, err)
	}
	ci.Manifest = &manifest

	return nil
}

func (ci *ContainerImage) GetImage(all bool) error {
	out, err := ci.DClient.ImagePull(
		context.Background(),
		ci.URL,
		types.ImagePullOptions{
			All:          all,
			RegistryAuth: ci.Auth,
		},
	)
	if err != nil {
		return fmt.Errorf("Error grabbing container image %s: %v", ci.URL, err)
	}

	defer out.Close()

	if ci.Debug {
		_, err := io.Copy(os.Stdout, out)
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", ci.URL, err)
		}
	} else {
		var b bytes.Buffer

		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		f.Write(b.Bytes())
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", ci.URL, err)
		}
	}

	return nil
}

func (ci *ContainerImage) RetagImage(srcImage, destImage string) error {
	err := ci.DClient.ImageTag(context.Background(), srcImage, destImage)
	if err != nil {
		return fmt.Errorf("Error re-tagging public container image from %s to %s: %v", srcImage, destImage, err)
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

func (ci *ContainerImage) CheckContainerImage() error {

	if err := ci.GetImage(false); err != nil {
		return fmt.Errorf("Error getting container image metadata: %w", err)
	}

	if err := ci.EnsureSourceImage(); err != nil {
		return fmt.Errorf("Error ensuring source image: %w", err)
	}

	if err := ci.GetMetadata(); err != nil {
		return fmt.Errorf("Error getting container image metadata: %w", err)
	}

	return nil
}
