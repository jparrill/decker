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
	coreReg "github.com/jparrill/decker/pkg/core/registry"
)

func (reg *Registry) Verify() error {
	var ok bool

	ps := NewVerifyPullSecret(reg.FilePath, false, reg.Debug)

	reg.PSData, ok = ps.Data.Auths[reg.URL]
	if !ok {
		check.Checker("Find registry in pull secret", fmt.Errorf("Registry %s not found in pull secret", reg.URL))
		return fmt.Errorf("Registry %s not found in pull secret", reg.URL)
	}

	err := reg.VerifyRegistryCredentials()
	check.Checker("Registry Authentication", err)
	if err != nil {
		return err
	}

	err = reg.VerifyRegistryPushAndPull()
	if err != nil {
		return err
	}

	return nil
}

// VerifyRegistryCredentials verifies the registry credentials for the given RegistryEntry.
// It creates a new Docker client and logs in to the registry using the provided credentials.
// Returns an error if the client cannot be created or if the login fails.
func (reg *Registry) VerifyRegistryCredentials() error {

	dCli, err := dockerclient.NewClientWithOpts()
	if err != nil {
		return err
	}

	err = reg.FillAuthCredentials()
	if err != nil {
		return err
	}

	authConf := &dockerregistrytype.AuthConfig{
		ServerAddress: reg.URL,
		Username:      reg.PSData.Username,
		Password:      reg.PSData.Password,
		Auth:          reg.PSData.Auth,
		Email:         reg.PSData.Email,
	}

	_, err = dCli.RegistryLogin(context.Background(), *authConf)
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
func (reg *Registry) VerifyRegistryPushAndPull() error {

	if err := reg.Encode(); err != nil {
		return err
	}

	dCli, err := dockerclient.NewClientWithOpts(
		dockerclient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}

	ref, err := prepareTemporaryImage(dCli, "", reg.URL, reg.FilePath)
	check.Checker("Prepare temporary image", err)

	// Push
	rc, err := dCli.ImagePush(context.Background(), ref.String(), types.ImagePushOptions{
		All:          true,
		RegistryAuth: reg.EAuth,
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
		return fmt.Errorf("Cannot push image to registry %s", ref.Registry)
	}

	// Delete local image
	if _, err = dCli.ImageRemove(context.Background(), ref.String(), types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}); err != nil {
		return fmt.Errorf("Error removing image %s from local machine: %v", ref.String(), err)
	}

	cImage := NewVerifyContainerImage(ref.String(), reg.EAuth, reg.FilePath, dCli)

	// Make sure the cImage is not in local
	// If err means that the image is not there
	if err := cImage.EnsureSourceImage(); err == nil {
		return err
	}

	// Pull
	err = cImage.GetImage(true)
	check.Checker("Registry Pull Permissions", err)

	return nil
}

func NewVerifyRegistry(url, filePath string, debug bool) (*Registry, error) {
	reg, err := coreReg.NewRegistry(url, filePath, false, debug)
	if err != nil {
		return nil, err
	}

	return &Registry{
		Registry: *reg,
	}, nil
}
