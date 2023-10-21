package verify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	dockerregistrytype "github.com/docker/docker/api/types/registry"
	dockerclient "github.com/docker/docker/client"

	"github.com/jparrill/decker/pkg/core/check"
)

const (
	alpineSampleImage = "quay.io/libpod/alpine:latest"
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

			if err := VerifyRegistryCredentials(registryName, record); err != nil {
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
		check.Checker("Find registry in pull secret", fmt.Errorf("registry %s not found in pull secret", o.Registry))
	}

	err = VerifyRegistryCredentials(o.Registry, registryData)
	check.Checker("Registry Authentication", err)

	return nil

}

func VerifyRegistryCredentials(registryURL string, record RegistryRecordType) error {

	dCli, err := dockerclient.NewClientWithOpts()
	if err != nil {
		return err
	}

	fillAuthCredentials(&record)

	authConfig := dockerregistrytype.AuthConfig{
		ServerAddress: registryURL,
		Username:      record.Username,
		Password:      record.Password,
		Auth:          record.Auth,
	}

	_, err = dCli.RegistryLogin(context.Background(), authConfig)
	if err != nil {
		return err
	}

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

func fillAuthCredentials(record *RegistryRecordType) error {

	if len(record.Username) <= 0 || len(record.Password) <= 0 {
		authBytes, err := base64.StdEncoding.DecodeString(record.Auth)
		if err != nil {
			return err
		}
		authPair := strings.Split(string(authBytes), ":")
		if len(authPair) != 2 {
			return fmt.Errorf("Bad formed authentication token")
		}
		record.Username = authPair[0]
		record.Password = authPair[1]
	}

	if len(record.Auth) <= 0 {
		authString := fmt.Sprintf("%s:%s", record.Username, record.Password)
		record.Auth = base64.StdEncoding.EncodeToString([]byte(authString))
	}

	return nil

}
