package ocprelease

import (
	"fmt"
	"log"

	"github.com/jparrill/decker/pkg/core/check"
	"github.com/jparrill/decker/pkg/validate"
	"github.com/spf13/cobra"
)

func NewValidateCommand() *cobra.Command {

	ocpImage := validate.OCPImage{}

	cmd := &cobra.Command{
		Use:          "ocp-release",
		Short:        "Validates an OCP release payload",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&ocpImage.URL, "url", "", "Registry url to check access to.")
	cmd.Flags().BoolVar(&ocpImage.Debug, "debug", false, "Debug mode to verify the Pull Secret")
	cmd.Flags().StringVar(&ocpImage.PullSecret.FilePath, "authfile", "", "Path to the pull secret file")
	if err := cmd.MarkFlagRequired("url"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Verifying Image: " + check.BoldWhite.Render(ocpImage.URL))
		if err := ocpImage.Validate(); err != nil {
			fmt.Println()
			return err
		}

		return nil
	}

	return cmd
}
