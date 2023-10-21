package pullsecret

import (
	"fmt"
	"log"

	verifypkg "github.com/jparrill/decker/pkg/verify"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {

	opts := verifypkg.PullSecretOpts{}

	cmd := &cobra.Command{
		Use:          "pull-secret",
		Short:        "Verifies a pullsecret structure",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&opts.File, "authfile", "", "Path to the pull secret file")
	cmd.Flags().BoolVar(&opts.Inspect, "inspect", false, "Check the registries details included in PullSecret file")
	if err := cmd.MarkFlagRequired("authfile"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Verifying pullsecret: %s\n", opts.File)
		verifypkg.VerifyPullSecret(opts)

		return nil
	}

	return cmd
}
