package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func getConfigPath() string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return filepath.Join(dir, ".fastlanerc")
}

//NewConfig is the factory for Config
func NewConfig() (*Config, error) {
	configPath := getConfigPath()
	var config *Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config = &Config{}
		err := config.Serialize()
		if err != nil {
			return nil, err
		}
	} else {
		content, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(content, &config)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

//Host - fastlane host
type Host struct {
	Name string
	URL  string
}

//Config loads or creates a configuration in the user's home
type Config struct {
	Hosts map[string]*Host
}

//List all the configured servers
func (config *Config) List() map[string]*Host {
	return config.Hosts
}

//AddTarget to .fastlanerc
func (config *Config) AddTarget(name, url string) {
	config.Hosts[name] = &Host{
		Name: name,
		URL:  url,
	}
}

//Serialize the config to the home directory
func (config *Config) Serialize() error {
	configPath := getConfigPath()
	var result []byte

	result, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, result, 0644)
	if err != nil {
		return err
	}

	return nil
}
