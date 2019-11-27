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
func (inv *Inventory) Put(labels map[string]string) error {
	for _, item := range inv.Items {
		if item["name"] == labels["name"] {
			for k, v := range labels {
				item[k] = v
			}
			return nil
		}
	}

	inv.Items = append(inv.Items, labels)
	return nil
}
