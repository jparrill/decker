package verify

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jparrill/decker/pkg/core/check"
)

func (ps *PullSecret) Verify() {
	var data AuthsType

	jsonData, err := os.ReadFile(ps.FilePath)
	check.Checker("Read input file", err)

	err = json.Unmarshal(jsonData, &data)
	check.Checker("Unmarshal JSON file", err)

	fmt.Println()

	if ps.Inspect {
		for registryName, record := range data.Auths {
			registryEntry := NewRegistryAuth(registryName, record.Username, record.Password, record.Auth)

			fmt.Println("RegistryName: " + check.BoldWhite.Render(registryName))
			if len(record.Auth) <= 0 && (len(record.Username) <= 0 && len(record.Password) <= 0) {
				check.Checker("Registry Credentials", fmt.Errorf("No authentication provided"))
			} else {
				check.Checker("Registry Credentials", nil)
			}

			if err := registryEntry.VerifyRegistryCredentials(); err != nil {
				check.Checker("Registry Authentication", fmt.Errorf("Error login into destination registry"))
			}

			fmt.Println()
		}
	}
}

func (ps *PullSecret) GetPullSecretData() (AuthsType, error) {
	var data AuthsType

	jsonData, err := os.ReadFile(ps.FilePath)
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
