package ocprelease

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewValidateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "ocp-release",
		Short:        "Validates an OCP release payload",
		SilenceUsage: true,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		//ctx := cmd.Context()
		//if opts.Timeout > 0 {
		//	var cancel context.CancelFunc
		//	ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		//	defer cancel()
		//}

		fmt.Println("validate ocprelease")
		return nil
	}

	return cmd
}
