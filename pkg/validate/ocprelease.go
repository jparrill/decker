package validate

import (
	"fmt"

	"github.com/jparrill/decker/pkg/core/check"
	coreReg "github.com/jparrill/decker/pkg/core/registry"
	"github.com/jparrill/decker/pkg/verify"
	"github.com/openshift/library-go/pkg/image/reference"
)

func (oi *OCPImage) Validate() error {
	var err error

	ps := verify.NewVerifyPullSecret(oi.FilePath, false, oi.Debug)

	oi.SrcRegistry = coreReg.NewRegistry(
		oi.SrcRegistry.URL,
		oi.SrcRegistry.FilePath,
		oi.SrcRegistry.TLSVerify,
		oi.SrcRegistry.Debug,
	)

	oi.SrcRegistry.PSData = ps.Data.Auths[oi.SrcRegistry.URL]

	if err := oi.SrcRegistry.Encode(); err != nil {
		return fmt.Errorf("Error marshalling auth credentials: %w", err)
	}

	ref, err := reference.Parse(oi.URL)
	check.Checker("Container Image Parsed", err)
	if err != nil {
		return fmt.Errorf("Error parsing container image: %w", err)
	}

	oi.Ref = &ref

	if _, ok := ps.Data.Auths[oi.SrcRegistry.URL]; !ok {
		return fmt.Errorf("Registry %s not found in pull secret", oi.SrcRegistry.URL)
	}
	check.Checker("Registry found in PullSecret", err)

	if err := oi.CheckContainerImage(
		oi.SrcRegistry,
	); err != nil {
		return fmt.Errorf("Error checking container image: %w", err)
	}
	check.Checker(fmt.Sprintf("Container Image validated %s", oi.Ref.Name), err)

	return nil
}

func (oi *OCPImage) CheckContainerImage(reg *coreReg.Registry) error {

	if err := oi.GetMetadata(); err != nil {
		return fmt.Errorf("Error getting container image metadata: %w", err)
	}

	fmt.Println("Manifest:", oi.Manifest)

	//srcDigest, err := getDigestFromRegistry(publicImageName, publicRegistryURL)
	//if err != nil {
	//	fmt.Println("Error al obtener el digest del registro público:", err)
	//	return
	//}

	//// Obtiene el digest del registro privado
	//destDigest, err := getDigestFromRegistry(privateImageName, privateRegistryURL)
	//if err != nil {
	//	fmt.Println("Error al obtener el digest del registro privado:", err)
	//	return
	//}

	return nil
}
