package backend

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type Inventory struct {
	Items     []Item            `yaml:"items"`
	Templates map[string]string `yaml:"templates"`
	Shortcuts map[string]string `yaml:"shortcuts"`
}

type Item map[string]string

func loadInventory(filename string) (*Inventory, error) {
	inv := &Inventory{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, inv)
	if err != nil {
		return nil, err
	}
	logrus.Tracef("inventory path: %s struct: %#v, ", filename, inv)
	return inv, nil
}
func saveInventory(inv *Inventory, filename string) error {
	buf, err := yaml.Marshal(inv)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, buf, 0644)

}
func (inv *Inventory) Init() *Inventory {
	if inv == nil {
		inv = &Inventory{}
	}

	if inv.Items == nil {
		inv.Items = []Item{}
	}

	if inv.Templates == nil {
		inv.Templates = map[string]string{}
	}

	if _, ok := inv.Templates["default"]; !ok {
		inv.Templates["default"] = "{{args}}"
	}

	return inv
}
