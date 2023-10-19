package generate

import (
	"github.com/jparrill/decker/cmd/generate/icsp"
	"github.com/jparrill/decker/cmd/generate/idms"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "generate",
		Short:        "Creates the ICSP/IDMS manifests for an OCP release mirrored.",
		SilenceUsage: true,
	}

	cmd.AddCommand(icsp.NewGenerateCommand())
	cmd.AddCommand(idms.NewGenerateCommand())

	return cmd
}
