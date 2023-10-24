package verify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func NewRegistryAuth(registryURL, user, pass, auth string) *RegistryEntry {
	return &RegistryEntry{
		ServerAddress: registryURL,
		Username:      user,
		Password:      pass,
		Auth:          auth,
	}
}

func (rge *RegistryEntry) FillAuthCredentials() error {

	if len(rge.Username) <= 0 || len(rge.Password) <= 0 {
		authBytes, err := base64.StdEncoding.DecodeString(rge.Auth)
		if err != nil {
			return err
		}
		authPair := strings.Split(string(authBytes), ":")
		if len(authPair) != 2 {
			return fmt.Errorf("Bad formed authentication token")
		}
		rge.Username = authPair[0]
		rge.Password = authPair[1]
	}

	if len(rge.Auth) <= 0 {
		authString := fmt.Sprintf("%s:%s", rge.Username, rge.Password)
		rge.Auth = base64.StdEncoding.EncodeToString([]byte(authString))
	}

	return nil
}

func (rge *RegistryEntry) Encode() (string, error) {

	encodedJSON, err := json.Marshal(rge)
	if err != nil {
		return "", fmt.Errorf("Error marshalling authconfig: %v", err)
	}

	return base64.URLEncoding.EncodeToString(encodedJSON), nil

}
