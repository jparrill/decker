package types

import (
	"net/url"
)

type RegistryOpts struct {
	Destination *url.URL
	PullSecret  []byte
	TLSVerify   bool
}

type PullSecretOpts struct {
	File string
}

type ImageOpts struct {
	ImageURL   *url.URL
	PullSecret []byte
	TLSVerify  bool
}
