package pullsecret

import (
	"fmt"
	"log"
	"os"

	"github.com/jparrill/decker/pkg/verify"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {

	pullSecret := verify.PullSecret{}

	cmd := &cobra.Command{
		Use:          "pull-secret",
		Short:        "Verifies a pullsecret structure",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&pullSecret.FilePath, "authfile", "", "Path to the pull secret file")
	cmd.Flags().BoolVar(&pullSecret.Inspect, "inspect", false, "Check the registries details included in PullSecret file")
	cmd.Flags().BoolVar(&pullSecret.Debug, "debug", false, "Debug mode to verify the Pull Secret")
	if err := cmd.MarkFlagRequired("authfile"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Verifying pullsecret: %s\n", pullSecret.FilePath)
		if err := pullSecret.Verify(); err != nil {
			if pullSecret.Debug {
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}

		return nil
	}

	return cmd
}
