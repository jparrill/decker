package validate

import (
	"github.com/jparrill/decker/cmd/validate/ocprelease"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "validate",
		Short:        "Validate an OCP payload images and their layers",
		SilenceUsage: true,
	}

	cmd.AddCommand(ocprelease.NewValidateCommand())

	return cmd
}
