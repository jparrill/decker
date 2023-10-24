package image

import (
	"fmt"
	"log"
	"os"

	"github.com/jparrill/decker/pkg/core/check"
	"github.com/jparrill/decker/pkg/verify"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {
	image := verify.ContainerImage{}

	cmd := &cobra.Command{
		Use:          "image",
		Short:        "Verifies if a container image is in a destination registry",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&image.URL, "url", "", "Registry url to check access to.")
	cmd.Flags().StringVar(&image.FilePath, "authfile", "", "Path to the pull secret file")
	cmd.Flags().BoolVar(&image.Debug, "debug", false, "Debug mode to verify the Pull Secret")
	if err := cmd.MarkFlagRequired("url"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Verifying Image: " + check.BoldWhite.Render(image.URL))
		if err := image.Verify(); err != nil {
			if image.Debug {
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}
		return nil
	}

	return cmd
}
