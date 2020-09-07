package backend

import (
	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
)

type ShortcutSetOptions struct {
	Name string
	Expr string
}

func (backend *Backend) ShortcutSet(opt *ShortcutSetOptions) error {
	inv, _ := loadInventory(backend.invFilePath)
	// ignore error
	inv = inv.Init()

	// TODO vaildate
	_, err := expr.Compile(opt.Expr)
	if err != nil {
		errors.Wrapf(err, "invalid expr")
	}
	inv.Shortcuts[opt.Name] = opt.Expr

	return saveInventory(inv, backend.invFilePath)
}

type ShortcutDelOptions struct {
	Names []string
}

func (backend *Backend) ShortcutDel(opt *ShortcutDelOptions) error {
	inv, _ := loadInventory(backend.invFilePath)
	// ignore error
	inv = inv.Init()

	for _, name := range opt.Names {
		delete(inv.Shortcuts, name)
	}

	return saveInventory(inv, backend.invFilePath)
}
