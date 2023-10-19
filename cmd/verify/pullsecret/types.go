package pullsecret

type AuthsType struct {
	Auths map[string]RegistryRecordType
}

type RegistryRecordType struct {
	Auth     string `json:"auth,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}
