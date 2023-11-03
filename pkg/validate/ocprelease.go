package validate

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/docker"
	referencev5 "github.com/containers/image/v5/docker/reference"

	//"github.com/containers/image/docker/reference"
	"github.com/containers/image/v5/types"
	"github.com/jparrill/decker/pkg/core/check"
	"github.com/jparrill/decker/pkg/verify"
	"github.com/openshift/library-go/pkg/image/reference"
)

func (oi *OCPImage) Validate() error {
	var ps verify.AuthsType
	var err error

	if ps, err = oi.PullSecret.GetPullSecretData(); err != nil {
		return err
	}
	check.Checker("Pull Secret", err)

	ref, err := reference.Parse(oi.URL)
	check.Checker("Container Image Parsed", err)
	if err != nil {
		return fmt.Errorf("Error parsing container image: %w", err)
	}

	if _, ok := ps.Auths[ref.Registry]; !ok {
		return fmt.Errorf("Registry %s not found in pull secret", ref.String())
	}
	check.Checker("Registry found in PullSecret", err)

	if err := oi.CheckOCPImage(ps.Auths[ref.Registry]); err != nil {
		return err
	}
	check.Checker("OCP Version validated", err)

	return nil
}

func (oi *OCPImage) CheckOCPImage(vrt verify.RegistryRecordType) error {

	registryEntry := verify.NewRegistryAuth(
		oi.URL,
		vrt.Username,
		vrt.Password,
		vrt.Auth,
	)
	err := registryEntry.FillAuthCredentials()
	if err != nil {
		return fmt.Errorf("Error filling auth credentials: %w", err)
	}

	sys := &types.SystemContext{
		DockerAuthConfig: &types.DockerAuthConfig{
			Username: registryEntry.Username,
			Password: registryEntry.Password,
		},
		DockerInsecureSkipTLSVerify: types.NewOptionalBool(true),
	}
	named, err := referencev5.ParseDockerRef(oi.URL)
	if err != nil {
		return fmt.Errorf("Error parsing dockerRef: %v\n", err)
	}
	srcRef, err := docker.NewReference(named)
	if err != nil {
		return fmt.Errorf("Error creating image reference: %v\n", err)
	}

	// Get the image manifest
	typesImageSource, err := srcRef.NewImageSource(context.Background(), sys)
	if err != nil {
		return fmt.Errorf("Error getting image manifest: %v\n", err)
	}

	data, name, err := typesImageSource.GetManifest(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error getting image manifest: %v\n", err)
	}

	fmt.Println("Manifest: ", string(data))
	fmt.Println("Name: ", name)

	return nil
}
