package registry

import (
	"fmt"
	"log"

	verifypkg "github.com/jparrill/decker/pkg/verify"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {

	opts := verifypkg.RegistryOpts{}

	cmd := &cobra.Command{
		Use:          "registry",
		Short:        "Verifies a registry accesibility",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&opts.Registry, "url", "", "Registry url to check access to.")
	cmd.Flags().StringVar(&opts.PullSecretOpts.File, "authfile", "", "Pull secret to authenticate against the destination registry")
	cmd.Flags().BoolVar(&opts.Insecure, "insecure", false, "Allow insecure registry connections.")
	if err := cmd.MarkFlagRequired("url"); err != nil {
		log.Fatal(err)
	}

	if err := cmd.MarkFlagRequired("authfile"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("verify registry")
		verifypkg.VerifyRegistry(opts)
		return nil
	}

	return cmd
}
