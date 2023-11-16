package validate

import (
	"fmt"

	"github.com/jparrill/decker/pkg/core/check"
)

func Validate(SrcImage, DstImage *OCPImage) error {

	err := SrcImage.CheckContainerImage()
	if err != nil {
		return fmt.Errorf("Error checking container image: %w", err)
	}
	check.Checker(fmt.Sprintf("Source Container Image validated %s", SrcImage.Ref.String()), err)

	err = DstImage.CheckContainerImage()
	if err != nil {
		return fmt.Errorf("Error checking container image: %w", err)
	}
	check.Checker(fmt.Sprintf("Destination Container Image validated %s", SrcImage.Ref.String()), err)

	if SrcImage.ContainerImage.Manifest.ID != DstImage.ContainerImage.Manifest.ID {
		check.Checker(
			"Source and Destination Container Images validated",
			fmt.Errorf("Source (%s) and Destination (%s) Container Images are different", SrcImage.ContainerImage.Manifest.ID, DstImage.ContainerImage.Manifest.ID),
		)
		return fmt.Errorf("Source (%s) and Destination (%s) Container Images are different", SrcImage.ContainerImage.Manifest.ID, DstImage.ContainerImage.Manifest.ID)
	}

	check.Checker("Source and Destination Container Images has the same digest", nil)
	return nil
}
