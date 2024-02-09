package diagnostic

type Config struct {
	ConfigFilePath string
	Content        *Content
}

type Content []byte
