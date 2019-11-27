package backend

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

func (backend *Backend) Get(expr string, outputFormat string) error {
	inv, err := loadInventory(backend.invFilePath)
	if err != nil {
		return err
	}
	inv = inv.Init()

	items, err := inv.ApplyExpr(expr)
	if err != nil {
		return err
	}

	switch outputFormat {
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
