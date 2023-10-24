package verify

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	dockerregistrytype "github.com/docker/docker/api/types/registry"
	dockerclient "github.com/docker/docker/client"

	"github.com/jparrill/decker/pkg/core/check"
)

const (
	alpineSampleImage = "quay.io/libpod/alpine:latest"
	debug             = false
)

func (rg *Registry) Verify() error {

	ps := &PullSecret{
		FilePath: rg.FilePath,
	}

	psData, err := ps.GetPullSecretData()
	if err != nil {
		return fmt.Errorf("Error getting pull secret data: %w", err)
	}

	registryData, ok := psData.Auths[rg.URL]
	if !ok {
		check.Checker("Find registry in pull secret", fmt.Errorf("Registry %s not found in pull secret", rg.URL))
	}

	registryEntry := NewRegistryAuth(
		rg.URL,
		registryData.Username,
		registryData.Password,
		registryData.Auth,
	)

	err = registryEntry.VerifyRegistryCredentials()
	check.Checker("Registry Authentication", err)

	err = registryEntry.VerifyRegistryPushAndPull()
	if err != nil {
		return err
	}

	return nil

}

// VerifyRegistryCredentials verifies the registry credentials for the given RegistryEntry.
// It creates a new Docker client and logs in to the registry using the provided credentials.
// Returns an error if the client cannot be created or if the login fails.
func (rge *RegistryEntry) VerifyRegistryCredentials() error {

	dCli, err := dockerclient.NewClientWithOpts()
	if err != nil {
		return err
	}

	err = rge.FillAuthCredentials()
	if err != nil {
		return err
	}

	_, err = dCli.RegistryLogin(context.Background(), dockerregistrytype.AuthConfig(*rge))
	if err != nil {
		return err
	}

	return nil
}

// VerifyRegistryPushAndPull verifies if the registry
// - can be pushed and pulled by encoding the private registry auth,
// - preparing a temporary image
// - pushing the image to the registry
// - pulling the image from the registry.
func (rge *RegistryEntry) VerifyRegistryPushAndPull() error {

	privateRegistryAuth, err := rge.Encode()
	if err != nil {
		return fmt.Errorf("failed creating auth for image push: %w", err)
	}

	dCli, err := dockerclient.NewClientWithOpts(
		dockerclient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}

	ref, err := prepareTemporaryImage(dCli, "", rge.ServerAddress)
	check.Checker("Prepare temporary image", err)

	// Push
	rc, err := dCli.ImagePush(context.Background(), ref.String(), types.ImagePushOptions{
		All:          true,
		RegistryAuth: privateRegistryAuth,
	})

	defer rc.Close()

	var out bytes.Buffer

	_, err = io.Copy(&out, rc)
	if err != nil {
		return fmt.Errorf("Error writting image to buffer %s: %v", ref.String(), err)
	}

	if debug {
		f := bufio.NewWriter(os.Stdout)
		defer f.Flush()
		_, err := f.Write(out.Bytes())
		if err != nil {
			return fmt.Errorf("Error writting image to buffer %s: %v", ref.String(), err)
		}
	}

	if strings.Contains(out.String(), "error") {
		check.Checker("Registry Push Permissions", fmt.Errorf("Cannot push image to registry %s", ref.Registry))
	}

	// Delete local image
	if _, err = dCli.ImageRemove(context.Background(), ref.String(), types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}); err != nil {
		return fmt.Errorf("Error removing image %s from local machine: %v", ref.String(), err)
	}

	cImage := &ContainerImage{
		DClient: dCli,
		URL:     ref.String(),
		Auth:    privateRegistryAuth,
	}

	// Make sure the cImage is not in local
	// If err means that the image is not there
	if err := cImage.EnsureSourceImage(); err == nil {
		return err
	}

	// Pull
	err = cImage.GetImage()
	check.Checker("Registry Pull Permissions", err)

	return nil
}
