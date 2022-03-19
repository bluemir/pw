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
	_, err := expr.Compile(opt.Expr)
	if err != nil {
		errors.Wrapf(err, "invalid expr")
	}
	backend.inv.Shortcuts[opt.Name] = opt.Expr

	return backend.Save()
}

type ShortcutDelOptions struct {
	Names []string
}

func (backend *Backend) ShortcutDelete(opt *ShortcutDelOptions) error {
	for _, name := range opt.Names {
		delete(backend.inv.Shortcuts, name)
	}

	return backend.Save()
}
