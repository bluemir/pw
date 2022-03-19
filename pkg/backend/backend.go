package backend

import (
	"github.com/sirupsen/logrus"
)

type Backend struct {
	invFilePath string
	inv         *Inventory
}

func New(invFilePath string) (*Backend, error) {
	inv, err := loadInventory(invFilePath)
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
	return backend.inv.save(backend.invFilePath)
}
