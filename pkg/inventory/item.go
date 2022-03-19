package inventory

import (
	"github.com/antonmedv/expr"
	"github.com/sirupsen/logrus"
)

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
