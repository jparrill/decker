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
		SRCImage string
		DSTImage string
		FilePath string
		Debug    bool
	)

	cmd := &cobra.Command{
		Use:          "ocp-release",
		Short:        "Validates an OCP release payload",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&SRCImage, "src", "", "Source container Image to validate")
	cmd.Flags().StringVar(&DSTImage, "dst", "", "Destination container Image to validate")
	cmd.Flags().BoolVar(&Debug, "debug", false, "Debug mode to verify the Pull Secret")
	cmd.Flags().StringVar(&FilePath, "authfile", "", "Path to the pull secret file")
	if err := cmd.MarkFlagRequired("src"); err != nil {
		log.Fatal(err)
	}

	if err := cmd.MarkFlagRequired("dst"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Validating Images: ")
		fmt.Printf("Source: " + check.BoldWhite.Render(SRCImage) + "\n")
		fmt.Printf("Destination: " + check.BoldWhite.Render(DSTImage) + "\n")

		srcImage := validate.NewValidateOCPImage(SRCImage, "", FilePath, nil)
		if Debug {
			srcImage.Debug = true
		}

		dstImage := validate.NewValidateOCPImage(DSTImage, "", FilePath, nil)
		if Debug {
			dstImage.Debug = true
		}

		if err := validate.Validate(srcImage, dstImage); err != nil {
			fmt.Println()
			return err
		}

		return nil
	}

	return cmd
}
