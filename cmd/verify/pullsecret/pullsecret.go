package pullsecret

import (
	"fmt"
	"log"

	"github.com/jparrill/decker/pkg/verify"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {

	var (
		FilePath string
		Inspect  bool
		Debug    bool
	)

	cmd := &cobra.Command{
		Use:          "pull-secret",
		Short:        "Verifies a pullsecret structure",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&FilePath, "authfile", "", "Path to the pull secret file")
	cmd.Flags().BoolVar(&Inspect, "inspect", false, "Check the registries details included in PullSecret file")
	cmd.Flags().BoolVar(&Debug, "debug", false, "Debug mode to verify the Pull Secret")
	if err := cmd.MarkFlagRequired("authfile"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Verifying pullsecret: %s\n", FilePath)
		pullsecret := verify.NewVerifyPullSecret(FilePath, Inspect, Debug)
		if err := pullsecret.Verify(); err != nil {
			if len(err) > 1 {
				return fmt.Errorf("Multiple errors found...")
			} else {
				return err[0]
			}
		}

		return nil
	}

	return cmd
}
