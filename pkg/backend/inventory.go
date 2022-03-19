package backend

import (
	"io/ioutil"

	"github.com/antonmedv/expr"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Inventory struct {
	Items     []Item            `yaml:"items"`
	Templates map[string]string `yaml:"templates"`
	Shortcuts map[string]string `yaml:"shortcuts"`
}

type Item map[string]string

// load inventory. if load failed return empty inventory
func loadInventory(filename string) (*Inventory, error) {
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
func (inv *Inventory) save(filename string) error {
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
		inv.Templates = map[string]string{}
	}

	if _, ok := inv.Templates["default"]; !ok {
		inv.Templates["default"] = "{{args}}"
	}

	return inv
}

func (inv *Inventory) GetItem(name string) Item {
	logrus.Trace(inv, name)
	for _, item := range inv.Items {
		if item["name"] == name {
			return item
		}
	}
	// not found
	return map[string]string{"name": name}
}

func (inv *Inventory) FindItem(e string) ([]Item, error) {
	if e == "" {
		e = "true"
	}
	program, err := expr.Compile(e)
	if err != nil {
		return nil, err
	}

	items := []Item{}
	for _, item := range inv.Items {
		output, err := expr.Run(program, item)
		if err != nil {
			return nil, err
		}

		if output.(bool) {
			items = append(items, item)
		}
	}

	logrus.Debug(items)

	return items, nil
}
func (inv *Inventory) SetItem(newItem Item) error {
	for _, item := range inv.Items {
		if item["name"] == newItem["name"] {
			for k, v := range newItem {
				item[k] = v
			}
			return nil
		}
	}

	inv.Items = append(inv.Items, newItem)

	return nil
}
func (inv *Inventory) DeleteItem(name string) {
	logrus.Trace(inv, name)
	for i, item := range inv.Items {
		if item["name"] == name {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			return
		}
	}
}
func (inv *Inventory) GetShortcut(name string) string {
	return inv.Shortcuts[name]
}
func (inv *Inventory) DeleteShortcut(name string) {
	delete(inv.Shortcuts, name)
}
func (inv *Inventory) SetShortcut(name string, expr string) {
	inv.Shortcuts[name] = expr
}
