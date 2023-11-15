package image

import (
	"fmt"
	"log"

	"github.com/jparrill/decker/pkg/core/check"
	"github.com/jparrill/decker/pkg/verify"

	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {

	var (
		URL      string
		FilePath string
		Debug    bool
	)

	cmd := &cobra.Command{
		Use:          "image",
		Short:        "Verifies if a container image is in a destination registry",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&URL, "url", "", "Registry url to check access to.")
	cmd.Flags().StringVar(&FilePath, "authfile", "", "Path to the pull secret file")
	cmd.Flags().BoolVar(&Debug, "debug", false, "Debug mode to verify the Pull Secret")
	if err := cmd.MarkFlagRequired("url"); err != nil {
		log.Fatal(err)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("Verifying Image: " + check.BoldWhite.Render(URL))
		image := verify.NewVerifyContainerImage(URL, "", FilePath, nil)
		if err := image.Verify(); err != nil {
			fmt.Println()
			return err
		}
		return nil
	}

	return cmd
}
