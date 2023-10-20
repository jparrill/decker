package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	dockerregistrytype "github.com/docker/docker/api/types/registry"
	dockerclient "github.com/docker/docker/client"
	"github.com/jparrill/decker/pkg/core/check"
)

func VerifyPullSecret(o PullSecretOpts) error {
	var data AuthsType

	jsonData, err := os.ReadFile(o.File)
	check.Checker("Read input file", err)

	err = json.Unmarshal(jsonData, &data)
	check.Checker("Unmarshal JSON file", err)

	fmt.Println()

	if o.DissectRegistry {
		for registryName, record := range data.Auths {
			fmt.Println("RegistryName: " + check.BoldWhite.Render(registryName))
			if len(record.Auth) <= 0 && (len(record.Username) <= 0 && len(record.Password) <= 0) {
				check.Checker("Registry Credentials", fmt.Errorf("No authentication provided"))
			} else {
				check.Checker("Registry Credentials", nil)
			}

			if err := VerifyRegistryCredentials(registryName, record); err != nil {
				check.Checker("Registry Authentication", fmt.Errorf("Error login into destination registry"))
			}

			fmt.Println()
		}
	}

	return nil
}

func VerifyRegistryCredentials(registryURL string, record RegistryRecordType) error {
	dCli, err := dockerclient.NewClientWithOpts()
	if err != nil {
		return err
	}

	authConfig := dockerregistrytype.AuthConfig{
		Username: record.Username,
		Password: record.Password,
		Auth:     record.Auth,
	}

	_, err = dCli.RegistryLogin(context.Background(), authConfig)
	if err != nil {
		return err
	}

	return nil
}
