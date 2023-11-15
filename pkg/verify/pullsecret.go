package verify

import (
	"fmt"

	"github.com/jparrill/decker/pkg/core/check"
	corePs "github.com/jparrill/decker/pkg/core/pullsecret"
)

func NewVerifyPullSecret(filePath string, inspect, debug bool) *PullSecret {
	ps := corePs.NewPullSecret(filePath, inspect, debug)

	return &PullSecret{
		PullSecret: *ps,
	}
}

func (ps *PullSecret) Verify() []error {
	var errs []error

	fmt.Println()

	if ps.Inspect {
		for registryName, record := range ps.Data.Auths {
			reg := NewVerifyRegistry(registryName, ps.FilePath, ps.Debug)
			reg.PSData = record
			if err := reg.Encode(); err != nil {
				check.Checker("Marshaling registry authentication", err)
				errs = append(errs, err)
			}

			fmt.Println("RegistryName: " + check.BoldWhite.Render(registryName))
			if len(record.Auth) <= 0 && (len(record.Username) <= 0 && len(record.Password) <= 0) {
				check.Checker("Registry Credentials", fmt.Errorf("No authentication provided"))
				errs = append(errs, fmt.Errorf("No authentication provided"))
			} else {
				check.Checker("Registry Credentials", nil)
			}

			if err := reg.VerifyRegistryCredentials(); err != nil {
				check.Checker("Registry Authentication", fmt.Errorf("Error login into destination registry"))
				errs = append(errs, fmt.Errorf("Error login into destination registry"))
			} else {
				check.Checker("Registry Authentication", nil)
			}

			fmt.Println()
		}
	}

	return errs
}
