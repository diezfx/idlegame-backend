package configloader

import (
	"fmt"
	"os"
)

type Source string

const (
	EnvSource  Source = "Env"
	FileSource Source = "File"
)

type Loader struct {
	Source
	configPath string
}

func NewFileLoader(configPath string) *Loader {
	return &Loader{configPath: configPath, Source: FileSource}
}

func (c *Loader) LoadConfig(namespace string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s.json", c.configPath, namespace)
	content, err := os.ReadFile(fmt.Sprintf("%s/%s.json", c.configPath, namespace))
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}
	return content, nil
}
