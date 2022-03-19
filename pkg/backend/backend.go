package backend

import (
	"github.com/sirupsen/logrus"

	"github.com/bluemir/pw/pkg/inventory"
)

type Backend struct {
	invFilePath string
	inv         *inventory.Inventory
}

func New(invFilePath string) (*Backend, error) {
	inv, err := inventory.Load(invFilePath)
	if err != nil {
		logrus.Warn("failed to load inventory")
		// skip.. continue with empty inventory
	}

	return &Backend{
		invFilePath,
		inv,
	}, err
}
func (backend *Backend) Save() error {
	return backend.inv.Save(backend.invFilePath)
}

type Item = inventory.Item
