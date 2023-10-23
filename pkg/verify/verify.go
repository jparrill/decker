package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"

	"github.com/jparrill/decker/pkg/core/check"
)

const (
	alpineSampleImage = "quay.io/libpod/alpine:latest"
	debug             = false
)

func VerifyPullSecret(o PullSecretOpts) {
	var data AuthsType

	jsonData, err := os.ReadFile(o.File)
	check.Checker("Read input file", err)

	err = json.Unmarshal(jsonData, &data)
	check.Checker("Unmarshal JSON file", err)

	fmt.Println()

	if o.Inspect {
		for registryName, record := range data.Auths {
			fmt.Println("RegistryName: " + check.BoldWhite.Render(registryName))
			if len(record.Auth) <= 0 && (len(record.Username) <= 0 && len(record.Password) <= 0) {
				check.Checker("Registry Credentials", fmt.Errorf("No authentication provided"))
			} else {
				check.Checker("Registry Credentials", nil)
			}

			if err := VerifyRegistryCredentials(registryName, &record); err != nil {
				check.Checker("Registry Authentication", fmt.Errorf("Error login into destination registry"))
			}

			fmt.Println()
		}
	}
}

func VerifyRegistry(o RegistryOpts) error {

	psData, err := GetPullSecretData(o.File)
	if err != nil {
		return fmt.Errorf("Error getting pull secret data: %w", err)
	}

	registryData, ok := psData.Auths[o.Registry]
	if !ok {
		check.Checker("Find registry in pull secret", fmt.Errorf("Registry %s not found in pull secret", o.Registry))
	}

	err = VerifyRegistryCredentials(o.Registry, &registryData)
	if err != nil {
		panic(err)
	}

	err = VerifyRegistryPushAndPull(o.Registry, &registryData)
	if err != nil {
		panic(err)
	}

	return nil

}

func VerifyRegistryCredentials(registryURL string, record *RegistryRecordType) error {

	dCli, err := dockerclient.NewClientWithOpts()
	if err != nil {
		return err
	}

	fillAuthCredentials(record)

	_, err = dCli.RegistryLogin(context.Background(), getRegistryAuth(registryURL, record))
	if err != nil {
		return err
	}
	check.Checker("Registry Authentication", err)

	return nil
}

func GetPullSecretData(authfile string) (AuthsType, error) {
	var data AuthsType

	jsonData, err := os.ReadFile(authfile)
	check.Checker("Read input file", err)
	if err != nil {
		return data, fmt.Errorf("Error reading input file: %v", err)
	}

	err = json.Unmarshal(jsonData, &data)
	check.Checker("Unmarshal JSON file", err)
	if err != nil {
		return data, fmt.Errorf("Error unmarshalling JSON file: %v", err)
	}

	return data, nil
}

func VerifyRegistryPushAndPull(registryURL string, record *RegistryRecordType) error {

	privateRegistryAuth, err := getEncodedRegistryAuth(registryURL, record)
	if err != nil {
		return fmt.Errorf("failed creating auth for image push: %w", err)
	}

	dCli, err := dockerclient.NewClientWithOpts(
		dockerclient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}

	ref, err := prepareTemporaryImage(dCli, "", registryURL)
	if err != nil {
		return fmt.Errorf("failed to parse image (%s): %w", alpineSampleImage, err)
	}

	// Push
	rc, err := dCli.ImagePush(context.Background(), ref.String(), types.ImagePushOptions{
		All:          true,
		RegistryAuth: privateRegistryAuth,
	})

	defer rc.Close()
	if debug {
		io.Copy(os.Stdout, rc)
	} else {
		var b bytes.Buffer
		io.Copy(&b, rc)
	}

	check.Checker("Registry Push Permissions", err)

	// Delete local image
	_, err = dCli.ImageRemove(context.Background(), ref.String(), types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})

	// Download image from custom registry again
	if err := ensureSourceImage(dCli, ref.String()); err != nil {
		return err
	}

	// Pull
	err = getImage(dCli, ref.String(), privateRegistryAuth)
	if err != nil {
		return err
	}

	check.Checker("Registry Pull Permissions", err)

	return nil
}
