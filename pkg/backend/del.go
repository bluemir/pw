package backend

func (backend *Backend) Del(itemNames []string) error {
	inv, _ := loadInventory(backend.invFilePath)
	inv = inv.Init()

	for _, name := range itemNames {
		inv.Del(name)
	}

	return saveInventory(inv, backend.invFilePath)
}
