package backend

import (
	"strings"
)

type TemplateSetOptions struct {
	Name string
	Args []string
}

func (backend *Backend) TemplateSet(opt *TemplateSetOptions) error {
	inv, _ := loadInventory(backend.invFilePath)
	// ignore error
	inv = inv.Init()

	inv.Templates[opt.Name] = strings.Join(opt.Args, " ")

	return saveInventory(inv, backend.invFilePath)
}

type TmeplateDelOptions struct {
	Names []string
}

func (backend *Backend) TemplateDel(opt *TemplateSetOptions) error {
	inv, _ := loadInventory(backend.invFilePath)
	// ignore error
	inv = inv.Init()

	delete(inv.Templates, opt.Name)

	return saveInventory(inv, backend.invFilePath)
}
