package failure

import (
	"errors"
	"time"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/log"
	yaml "gopkg.in/yaml.v2"
)

// AttackMap is a type that defines a list of map of attackers.
type AttackMap map[string]attack.Opts

// Config is the way a failure is defined.
type Config struct {
	Timeout time.Duration `yaml:"timeout,omitempty"`
	// used an array so the no repeated elements of map limitation can be bypassed.
	Attacks []AttackMap `yaml:"attacks,omitempty"`

	// TODO: accuracy
}

// ReadConfig Reads a config yaml defition and returns a config object.
func ReadConfig(data []byte) (Config, error) {
	log.Debug("reading config")
	c := &Config{}
	err := yaml.Unmarshal(data, c)
	return *c, err
}

// Render renders a yaml form a Config object.
func (c *Config) Render() ([]byte, error) {
	log.Debug("rendering config")

	// Check if there are more then one elements on the maps of the list.
	for _, a := range c.Attacks {
		if len(a) != 1 {
			return nil, errors.New("each attack map of the attack list needs to be a single map")
		}
	}
	// Marshal to yaml
	return yaml.Marshal(c)
}

// UnmarshalYAML wraps yaml lib unmarshalling to have extra validations.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Made to bypass the unmarshaling recursion.
	type plain Config
	cc := &Config{}
	if err := unmarshal((*plain)(cc)); err != nil {
		return err
	}

	// Check if there are more then one elements on the maps of the list.
	for _, a := range cc.Attacks {
		if len(a) != 1 {
			return errors.New("attacks format error, tip: check identantion and '-' indicator")
		}
	}

	*c = *cc
	return nil
}
