package pullsecret

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	verifytypes "github.com/jparrill/decker/cmd/verify/types"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {

	opts := verifytypes.PullSecretOpts{}

	cmd := &cobra.Command{
		Use:          "pull-secret",
		Short:        "Verifies a pullsecret structure",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&opts.File, "authfile", "", "Path to the pull secret file")
	err := cmd.MarkFlagRequired("authfile")
	if err != nil {
		log.Fatal(err)
	}

	// CONTINUE HERE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Verifying pullsecret: %s\n", opts.File)

		_ = verifyPullSecret(opts)

		return nil
	}

	return cmd
}

func verifyPullSecret(o verifytypes.PullSecretOpts) error {

	var data AuthsType

	jsonData, err := os.ReadFile(o.File)
	if err != nil {
		log.Fatalf("Error reading input file %s: %v\n", o.File, err)
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatalf("Error unmarshalling input file: %v", err)
	}

	for registryName, record := range data.Auths {
		fmt.Printf("RegistryName: %s\n", registryName)
		fmt.Printf("Auth: %s\n", record.Auth)
		fmt.Printf("Username: %s\n", record.Username)
		fmt.Printf("Password: %s\n", record.Password)
		fmt.Printf("Email: %s\n", record.Email)
		fmt.Println()
	}

	return nil
}
