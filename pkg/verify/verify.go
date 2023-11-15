package verify

import (
	coreImage "github.com/jparrill/decker/pkg/core/image"
	corePullSecret "github.com/jparrill/decker/pkg/core/pullsecret"
	coreReg "github.com/jparrill/decker/pkg/core/registry"
)

type Registry struct {
	coreReg.Registry
}

type ContainerImage struct {
	coreImage.ContainerImage
}

type PullSecret struct {
	corePullSecret.PullSecret
}
