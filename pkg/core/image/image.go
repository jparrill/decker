package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/jparrill/decker/pkg/core/check"
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

func NewContainerImage(url, auth, filepath string) (*ContainerImage, error) {
	ci := &ContainerImage{
		DClient:   nil,
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

func (ci *ContainerImage) GetMetadata() error {
	var err error
	var manifest types.ImageInspect

	ci.DClient, err = dockerclient.NewClientWithOpts()
	check.Checker("Docker Client Generated", err)
	if err != nil {
		return fmt.Errorf("Error generating docker client: %w", err)
	}

	manifest, _, err = ci.DClient.ImageInspectWithRaw(
		context.Background(),
		fmt.Sprintf("%s:%s", ci.Ref.Registry, ci.Ref.Name),
	)
	if err != nil {
		return fmt.Errorf("Error inspecting docker image %s: %w", fmt.Sprintf("%s:%s", ci.Ref.Registry, ci.Ref.Name), err)
	}
	ci.Manifest = &manifest

	return nil
}

func (ci *ContainerImage) GetImage() error {
	encodedJSON, err := json.Marshal(ci.Auth)
	if err != nil {
		return fmt.Errorf("Error marshalling authconfig: %v", err)
	}

	sDec := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := ci.DClient.ImagePull(
		context.Background(),
		ci.URL,
		types.ImagePullOptions{
			All:          true,
			RegistryAuth: sDec,
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
		_, err := io.Copy(&b, out)
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", ci.URL, err)
		}
	}

	return nil
}
