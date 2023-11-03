package validate

import (
	imageapi "github.com/openshift/api/image/v1"
)

type ReleaseImage struct {
	*imageapi.ImageStream `json:",inline"`
	StreamMetadata        *CoreOSStreamMetadata `json:"streamMetadata"`
}

type CoreOSStreamMetadata struct {
	Stream        string                        `json:"stream"`
	Architectures map[string]CoreOSArchitecture `json:"architectures"`
}

type CoreOSArchitecture struct {
	// Artifacts is a map of platform name to Artifacts
	Artifacts map[string]CoreOSArtifact `json:"artifacts"`
	Images    CoreOSImages              `json:"images"`
}

type CoreOSArtifact struct {
	Release string                             `json:"release"`
	Formats map[string]map[string]CoreOSFormat `json:"formats"`
}

type CoreOSFormat struct {
	Location           string `json:"location"`
	Signature          string `json:"signature"`
	SHA256             string `json:"sha256"`
	UncompressedSHA256 string `json:"uncompressed-sha256"`
}

type CoreOSImages struct {
	AWS      CoreOSAWSImages      `json:"aws"`
	PowerVS  CoreOSPowerVSImages  `json:"powervs"`
	Kubevirt CoreOSKubevirtImages `json:"kubevirt"`
}

type CoreOSAWSImages struct {
	Regions map[string]CoreOSAWSImage `json:"regions"`
}

type CoreOSAWSImage struct {
	Release string `json:"release"`
	Image   string `json:"image"`
}

type CoreOSKubevirtImages struct {
	Release   string `json:"release"`
	Image     string `json:"image"`
	DigestRef string `json:"digest-ref"`
}

type CoreOSPowerVSImages struct {
	Regions map[string]CoreOSPowerVSImage `json:"regions"`
}

type CoreOSPowerVSImage struct {
	Release string `json:"release"`
	Object  string `json:"object"`
	Bucket  string `json:"bucket"`
	URL     string `json:"url"`
}
