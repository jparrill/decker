package registry

import (
	"fmt"
	"log"

	"github.com/jparrill/decker/pkg/core/check"
	verify "github.com/jparrill/decker/pkg/verify"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {
	var (
		URL       string
		FilePath  string
		TLSVerify bool
		Debug     bool
	)

	cmd := &cobra.Command{
		Use:          "registry",
		Short:        "Verifies a registry accesibility",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&URL, "url", "", "Registry url to check access to.")
	cmd.Flags().StringVar(&FilePath, "authfile", "", "Pull secret to authenticate against the destination registry")
	cmd.Flags().BoolVar(&TLSVerify, "tls-verify", false, "Allow insecure registry connections.")
	cmd.Flags().BoolVar(&Debug, "debug", false, "Debug mode to verify the Pull Secret")
	if err := cmd.MarkFlagRequired("url"); err != nil {
		log.Fatal(err)
	}

	if err := cmd.MarkFlagRequired("authfile"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Verifying Registry: " + check.BoldWhite.Render(URL))
		registry := verify.NewVerifyRegistry(URL, FilePath, Debug)
		if err := registry.Verify(); err != nil {
			fmt.Println()
			return err
		}
		return nil
	}

	return cmd
}
