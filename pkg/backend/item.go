package backend

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type GetOptions struct {
	Expr         string
	ShortcutName string
	OutputFormat string
}

func (backend *Backend) GetItems(opt *GetOptions) error {
	if opt.ShortcutName != "" {
		if v, ok := backend.inv.Shortcuts[opt.ShortcutName]; ok {
			opt.Expr = v
		}
	}

	items, err := backend.inv.ApplyExpr(opt.Expr)
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
		logrus.Fatal("output is one of (name, yaml)")
	}

	return nil
}

type SetOptions struct {
	ItemNames []string
	Labels    map[string]string
}

func (backend *Backend) SetItems(opt *SetOptions) error {
	for _, name := range opt.ItemNames {
		item := backend.inv.Take(name)

		for k, v := range opt.Labels {
			if v == "-" {
				delete(item, k)
			} else {
				item[k] = v
			}
		}

		if err := backend.inv.Put(item); err != nil {
			return err
		}
	}

	return backend.Save()
}

type DeleteOptions struct {
	ItemNames []string
}

func (backend *Backend) DeleteItems(opt *DeleteOptions) error {
	for _, name := range opt.ItemNames {
		backend.inv.Delete(name)
	}
	return backend.Save()
}
