package registry

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "registry",
		Short:        "Verifies a registry accesibility",
		SilenceUsage: true,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		//ctx := cmd.Context()
		//if opts.Timeout > 0 {
		//	var cancel context.CancelFunc
		//	ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		//	defer cancel()
		//}

		fmt.Println("verify registry")
		return nil
	}

	return cmd
}
