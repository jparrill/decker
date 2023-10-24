package registry

import (
	"fmt"
	"log"

	"github.com/jparrill/decker/pkg/core/check"
	verifypkg "github.com/jparrill/decker/pkg/verify"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {

	registry := verifypkg.Registry{}

	cmd := &cobra.Command{
		Use:          "registry",
		Short:        "Verifies a registry accesibility",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&registry.URL, "url", "", "Registry url to check access to.")
	cmd.Flags().StringVar(&registry.FilePath, "authfile", "", "Pull secret to authenticate against the destination registry")
	cmd.Flags().BoolVar(&registry.Insecure, "insecure", false, "Allow insecure registry connections.")
	cmd.Flags().BoolVar(&registry.Debug, "debug", false, "Debug mode to verify the Pull Secret")
	if err := cmd.MarkFlagRequired("url"); err != nil {
		log.Fatal(err)
	}

	if err := cmd.MarkFlagRequired("authfile"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Verifying Registry: " + check.BoldWhite.Render(registry.URL))
		registry.Verify()
		return nil
	}

	return cmd
}
