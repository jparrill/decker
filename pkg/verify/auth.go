package verify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	dockerregistrytype "github.com/docker/docker/api/types/registry"
)

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

func getRegistryAuth(registryURL string, record *RegistryRecordType) dockerregistrytype.AuthConfig {
	return dockerregistrytype.AuthConfig{
		ServerAddress: registryURL,
		Username:      record.Username,
		Password:      record.Password,
		Auth:          record.Auth,
	}
}

func getEncodedRegistryAuth(registryURL string, record *RegistryRecordType) (string, error) {
	authConfig := getRegistryAuth(registryURL, record)

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", fmt.Errorf("Error marshalling authconfig: %v", err)
	}

	return base64.URLEncoding.EncodeToString(encodedJSON), nil

}
