package verify

import "net/url"

type AuthsType struct {
	Auths map[string]RegistryRecordType
}

type RegistryRecordType struct {
	Auth     string `json:"auth,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type RegistryOpts struct {
	Registry string
	Insecure bool
	PullSecretOpts
}

type PullSecretOpts struct {
	File    string
	Inspect bool
}

type ImageOpts struct {
	ImageURL   *url.URL
	PullSecret []byte
	TLSVerify  bool
}
