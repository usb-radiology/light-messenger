package configuration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

var config *Configuration

// Configuration ...
type Configuration struct {
	Server struct {
		HTTPPort int
	}
}

// LoadAndSetConfiguration ...
func LoadAndSetConfiguration(path string) (*Configuration, error) {
	var data Configuration

	file, readFileErr := ioutil.ReadFile(path)
	if readFileErr != nil {
		return nil, errors.Wrap(readFileErr, "could not read configuration file")
	}

	jsonUnmarshallErr := json.Unmarshal([]byte(file), &data)
	if jsonUnmarshallErr != nil {
		return nil, errors.Wrap(readFileErr, "could not parse json into config format")
	}
	config = &data
	return config, nil
}

// GetConfiguration ...
func GetConfiguration() (*Configuration, error) {
	if config == nil {
		return nil, errors.New("configuration empty / not loaded")
	}
	return config, nil
}
