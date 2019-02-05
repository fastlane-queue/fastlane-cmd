package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

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
		config = &Config{
			Hosts:           map[string]*Host{},
			LastUpdateCheck: 0,
		}
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
	Name     string
	URL      string
	Default  bool
	LastUsed int64
}

//Config loads or creates a configuration in the user's home
type Config struct {
	Hosts           map[string]*Host
	DefaultTarget   string
	LastUpdateCheck int64
}

//List all the configured servers
func (config *Config) List() map[string]*Host {
	return config.Hosts
}

//SetTarget to .fastlanerc
func (config *Config) SetTarget(name, url string, defaultTarget bool) {
	if defaultTarget {
		config.ClearDefaults()
		config.DefaultTarget = name
	}

	if _, ok := config.Hosts[name]; ok {
		config.Hosts[name].URL = url
		config.Hosts[name].Default = defaultTarget
	} else {
		config.Hosts[name] = &Host{
			Name:     name,
			URL:      url,
			Default:  defaultTarget,
			LastUsed: 0,
		}
	}
}

//ClearDefaults of targets
func (config *Config) ClearDefaults() {
	for _, host := range config.Hosts {
		host.Default = false
	}
	config.DefaultTarget = ""
}

// UpdateLastUsed for a given target
func (config *Config) UpdateLastUsed(name string) error {
	config.Hosts[name].LastUsed = time.Now().Unix()
	err := config.Serialize()
	if err != nil {
		return err
	}
	return nil
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
