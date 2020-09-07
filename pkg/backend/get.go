package backend

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type GetOptions struct {
	Expr             string
	ShortcutExprName string
	OutputFormat     string
}

func (backend *Backend) Get(opt *GetOptions) error {
	//func (backend *Backend) Get(expr string, shortcut string, outputFormat string) error {
	inv, err := loadInventory(backend.invFilePath)
	if err != nil {
		return err
	}
	inv = inv.Init()

	if opt.ShortcutExprName != "" {
		if v, ok := inv.Init().Shortcuts[opt.ShortcutExprName]; ok {
			opt.Expr = v
		}
	}

	items, err := inv.ApplyExpr(opt.Expr)
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
