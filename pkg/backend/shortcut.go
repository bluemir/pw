package backend

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/pkg/errors"
)

type GetShortcutOptions struct {
	Names []string
}

func (backend *Backend) GetShortcut(conf *GetShortcutOptions) error {
	if len(conf.Names) == 0 {
		scs := backend.inv.GetShortcuts()
		for k, v := range scs {
			fmt.Printf("%s: %s", k, v)
		}
		return nil
	}

	for _, name := range conf.Names {
		sc := backend.inv.GetShortcut(name)
		fmt.Printf("%s: %s", name, sc)
	}

	return nil
}

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
