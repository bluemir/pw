package backend

import (
	"github.com/antonmedv/expr"
	"github.com/sirupsen/logrus"
)

func (inv *Inventory) ApplyExpr(e string) ([]Item, error) {
	if e == "" {
		e = "true"
	}
	program, err := expr.Compile(e)
	if err != nil {
		return nil, err
	}

	items := []Item{}
	for _, item := range inv.Items {
		output, err := expr.Run(program, item)
		if err != nil {
			return nil, err
		}

		if output.(bool) {
			items = append(items, item)
		}
	}

	logrus.Debug(items)

	return items, nil
}
