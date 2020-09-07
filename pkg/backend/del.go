package backend

type DelOptions struct {
	ItemNames []string
}

func (backend *Backend) Del(opt *DelOptions) error {
	inv, _ := loadInventory(backend.invFilePath)
	inv = inv.Init()

	for _, name := range opt.ItemNames {
		inv.Del(name)
	}

	return saveInventory(inv, backend.invFilePath)
}
