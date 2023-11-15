package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type Registry struct {
	URL       string
	TLSVerify bool
	FilePath  string
	PSData    RegistryRecordType
	Debug     bool
	EAuth     string
}

// RegistryRecordType is the struct representing the PullSecretcomponents
type RegistryRecordType struct {
	Auth     string `json:"auth,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type RegistryEntryInterface interface {
	FillAuthCredentials() error
	Encode() error
}

func NewRegistry(url, filePath string, TLSVerify, debug bool) *Registry {
	return &Registry{
		URL:       url,
		TLSVerify: TLSVerify,
		FilePath:  filePath,
		PSData:    RegistryRecordType{},
		Debug:     debug,
	}
}

func (reg *Registry) FillAuthCredentials() error {

	if len(reg.PSData.Username) <= 0 || len(reg.PSData.Password) <= 0 {
		authBytes, err := base64.StdEncoding.DecodeString(reg.PSData.Auth)
		if err != nil {
			return err
		}
		authPair := strings.Split(string(authBytes), ":")
		if len(authPair) != 2 {
			return fmt.Errorf("Bad formed authentication token")
		}
		reg.PSData.Username = authPair[0]
		reg.PSData.Password = authPair[1]
	}

	if len(reg.PSData.Auth) <= 0 {
		authString := fmt.Sprintf("%s:%s", reg.PSData.Username, reg.PSData.Password)
		reg.PSData.Auth = base64.StdEncoding.EncodeToString([]byte(authString))
	}

	return nil
}

func (reg *Registry) Encode() error {

	encodedJSON, err := json.Marshal(reg.PSData)
	if err != nil {
		return fmt.Errorf("Error marshalling authconfig: %v", err)
	}

	reg.EAuth = base64.URLEncoding.EncodeToString(encodedJSON)

	return nil
}
