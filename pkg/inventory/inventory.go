package inventory

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Inventory struct {
	Items     []Item             `yaml:"items"`
	Templates map[string]Command `yaml:"templates"`
	Shortcuts map[string]string  `yaml:"shortcuts"`
}

type Item map[string]string
type Command []string

// load inventory. if load failed return empty inventory
func Load(filename string) (*Inventory, error) {
	inv := &Inventory{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Debug("inventory file not found")
		return inv.init(), err
	}

	if err = yaml.Unmarshal(content, inv); err != nil {
		logrus.Debug("inventory file unmarshal failed")
		return inv.init(), err
	}
	logrus.Tracef("inventory path: %s struct: %#v, ", filename, inv)
	return inv.init(), nil
}
func (inv *Inventory) Save(filename string) error {
	buf, err := yaml.Marshal(inv)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buf, 0644)
}
func (inv *Inventory) init() *Inventory {
	if inv == nil {
		inv = &Inventory{}
	}

	if inv.Items == nil {
		inv.Items = []Item{}
	}

	if inv.Templates == nil {
		inv.Templates = map[string]Command{}
	}

	if _, ok := inv.Templates["default"]; !ok {
		inv.Templates["default"] = []string{"{{args}}"}
	}

	return inv
}
