package backend

import "strings"

func (backend *Backend) Template(name string, args []string) error {
	inv, _ := loadInventory(backend.invFilePath)
	// ignore error
	inv = inv.Init()

	inv.Templates[name] = strings.Join(args, " ")

	return saveInventory(inv, backend.invFilePath)
}
