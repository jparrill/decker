package ocprelease

import (
	"fmt"
	"log"

	"github.com/jparrill/decker/pkg/core/check"
	"github.com/jparrill/decker/pkg/validate"
	"github.com/spf13/cobra"
)

func NewValidateCommand() *cobra.Command {

	var (
		URL      string
		FilePath string
		Debug    bool
	)

	cmd := &cobra.Command{
		Use:          "ocp-release",
		Short:        "Validates an OCP release payload",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&URL, "url", "", "Registry url to check access to.")
	cmd.Flags().BoolVar(&Debug, "debug", false, "Debug mode to verify the Pull Secret")
	cmd.Flags().StringVar(&FilePath, "authfile", "", "Path to the pull secret file")
	if err := cmd.MarkFlagRequired("url"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Validating Image: " + check.BoldWhite.Render(URL))
		ocpImage := validate.NewValidateOCPImage(URL, "", FilePath)

		if err := ocpImage.Validate(); err != nil {
			fmt.Println()
			return err
		}

		return nil
	}

	return cmd
}
