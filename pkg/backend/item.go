package backend

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type GetItemsOptions struct {
	Expr         string
	ShortcutName string
	OutputFormat string
}

func (backend *Backend) GetItems(opt *GetItemsOptions) error {
	if opt.ShortcutName != "" {
		if v, ok := backend.inv.Shortcuts[opt.ShortcutName]; ok {
			opt.Expr = v
		}
	}

	items, err := backend.inv.FindItem(opt.Expr)
	if err != nil {
		return err
	}

	switch opt.OutputFormat {
	case "name":
		for _, item := range items {
			fmt.Println(item["name"])
		}
	case "yaml":
		buf, err := yaml.Marshal(items)
		if err != nil {
			return err
		}
		fmt.Println(string(buf))
	case "json":
		buf, err := json.Marshal(items)
		if err != nil {
			return err
		}
		fmt.Println(string(buf))
	default:
		logrus.Fatal("output is one of (name, yaml, json)")
	}

	return nil
}

type SetItemOptions struct {
	ItemNames []string
	Labels    map[string]string
}

func (backend *Backend) SetItems(opt *SetItemOptions) error {
	for _, name := range opt.ItemNames {
		item := backend.inv.GetItem(name)

		for k, v := range opt.Labels {
			if v == "-" {
				delete(item, k)
			} else {
				item[k] = v
			}
		}

		if err := backend.inv.SetItem(item); err != nil {
			return err
		}
	}

	return backend.Save()
}

type DeleteItemOptions struct {
	ItemNames []string
}

func (backend *Backend) DeleteItems(opt *DeleteItemOptions) error {
	for _, name := range opt.ItemNames {
		backend.inv.DeleteItem(name)
	}
	return backend.Save()
}
