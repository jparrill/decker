package verify

import (
	dockerregistrytype "github.com/docker/docker/api/types/registry"
	dockerclient "github.com/docker/docker/client"
)

type AuthsType struct {
	Auths map[string]RegistryRecordType
}

// RegistryRecordType is the struct representing the PullSecretcomponents
type RegistryRecordType struct {
	Auth     string `json:"auth,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type Registry struct {
	URL      string
	Insecure bool
	FilePath string
	PSData   RegistryRecordType
	Debug    bool
}

type ContainerImage struct {
	DClient   *dockerclient.Client
	URL       string
	FilePath  string
	Auth      string
	TLSVerify bool
	Debug     bool
}

type PullSecret struct {
	FilePath string
	Inspect  bool
	Debug    bool
}

type Verifier interface {
	Verify()
}

type RegistryEntry dockerregistrytype.AuthConfig
