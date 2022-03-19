package backend

import (
	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
)

type SetShortcutOptions struct {
	Name string
	Expr string
}

func (backend *Backend) SetShortcut(opt *SetShortcutOptions) error {
	_, err := expr.Compile(opt.Expr)
	if err != nil {
		errors.Wrapf(err, "invalid expr")
	}
	backend.inv.SetShortcut(opt.Name, opt.Expr)

	return backend.Save()
}

type DeleteShortcutOptions struct {
	Names []string
}

func (backend *Backend) DeleteShortcut(opt *DeleteShortcutOptions) error {
	for _, name := range opt.Names {
		backend.inv.DeleteShortcut(name)
	}

	return backend.Save()
}
