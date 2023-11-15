package pullsecret

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jparrill/decker/pkg/core/check"
	coreReg "github.com/jparrill/decker/pkg/core/registry"
)

type AuthsType struct {
	Auths map[string]coreReg.RegistryRecordType
}

type PullSecret struct {
	FilePath string
	Inspect  bool
	Debug    bool
	Data     *AuthsType
}

type PullSecretInterface interface {
	GetPullSecretData()
	Encode(regName string) (string, error)
}

func NewPullSecret(filePath string, inspect, debug bool) *PullSecret {
	ps := &PullSecret{
		FilePath: filePath,
		Inspect:  inspect,
		Debug:    debug,
		Data:     nil,
	}

	ps.GetPullSecretData()

	return ps
}

func (ps *PullSecret) GetPullSecretData() {
	var data AuthsType

	jsonData, err := os.ReadFile(ps.FilePath)
	check.Checker("Read input file", err)
	if err != nil {
		panic(fmt.Errorf("Error reading input file: %v", err))
	}

	err = json.Unmarshal(jsonData, &data)
	check.Checker("Unmarshal JSON file", err)
	if err != nil {
		panic(fmt.Errorf("Error unmarshalling JSON file: %v", err))
	}

	ps.Data = &data
}

func (ps *PullSecret) Encode(regName string) (string, error) {

	encodedJSON, err := json.Marshal(ps.Data.Auths[regName])
	if err != nil {
		return "", fmt.Errorf("Error marshalling authconfig: %v", err)
	}

	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}
