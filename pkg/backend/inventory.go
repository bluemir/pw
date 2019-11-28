package backend

import "github.com/sirupsen/logrus"

func (inv *Inventory) Take(name string) Item {
	logrus.Trace(inv, name)
	for _, item := range inv.Items {
		if item["name"] == name {
			return item
		}
	}
	// not found
	return map[string]string{"name": name}
}
func (inv *Inventory) Put(newItem Item) error {
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
