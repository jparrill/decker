package icsp

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "icsp",
		Short:        "Generates the ICSP manifests",
		SilenceUsage: true,
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		//ctx := cmd.Context()
		//if opts.Timeout > 0 {
		//	var cancel context.CancelFunc
		//	ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		//	defer cancel()
		//}

		fmt.Println("generate ICSP")
		return nil
	}

	return cmd
}
