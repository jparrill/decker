package verify

import (
	"github.com/jparrill/decker/cmd/verify/image"
	"github.com/jparrill/decker/cmd/verify/pullsecret"
	"github.com/jparrill/decker/cmd/verify/registry"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "verify",
		Short:        "Verifies if an object is accesible, resilient or correct anatomivally talking",
		SilenceUsage: true,
	}

	cmd.AddCommand(pullsecret.NewVerifyCommand())
	cmd.AddCommand(registry.NewVerifyCommand())
	cmd.AddCommand(image.NewVerifyCommand())

	return cmd
}
