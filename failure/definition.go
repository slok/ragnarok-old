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

// Definition is the way a failure is defined.
type Definition struct {
	Timeout time.Duration `yaml:"timeout,omitempty"`
	// used an array so the no repeated elements of map limitation can be bypassed.
	Attacks []AttackMap `yaml:"attacks,omitempty"`

	// TODO: accuracy
}

// ReadDefinition Reads a config yaml defition and returns a definition object.
func ReadDefinition(data []byte) (Definition, error) {
	log.Debug("reading config")
	d := &Definition{}
	err := yaml.Unmarshal(data, d)
	return *d, err
}

// Render renders a yaml form a Definition object.
func (d *Definition) Render() ([]byte, error) {
	log.Debug("rendering config")

	// Check if there are more then one elements on the maps of the list.
	for _, a := range d.Attacks {
		if len(a) != 1 {
			return nil, errors.New("each attack map of the attack list needs to be a single map")
		}
	}
	// Marshal to yaml
	return yaml.Marshal(d)
}

// UnmarshalYAML wraps yaml lib unmarshalling to have extra validations.
func (d *Definition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Made to bypass the unmarshaling recursion.
	type plain Definition
	dd := &Definition{}
	if err := unmarshal((*plain)(dd)); err != nil {
		return err
	}

	// Check if there are more then one elements on the maps of the list.
	for _, a := range dd.Attacks {
		if len(a) != 1 {
			return errors.New("attacks format error, tip: check identantion and '-' indicator")
		}
	}

	*d = *dd
	return nil
}
