package idms

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "idms",
		Short:        "Generates the IDMS manifests",
		SilenceUsage: true,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		//ctx := cmd.Context()
		//if opts.Timeout > 0 {
		//	var cancel context.CancelFunc
		//	ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		//	defer cancel()
		//}

		fmt.Println("generate IDMS")
		return nil
	}

	return cmd
}
