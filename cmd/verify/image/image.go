package image

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "image",
		Short:        "Verifies if a container image is in a destination registry",
		SilenceUsage: true,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		//ctx := cmd.Context()
		//if opts.Timeout > 0 {
		//	var cancel context.CancelFunc
		//	ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		//	defer cancel()
		//}

		fmt.Println("verify image")
		return nil
	}

	return cmd
}
