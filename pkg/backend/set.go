package backend

func (backend *Backend) Set(itemNames []string, labels map[string]string) error {
	inv, _ := loadInventory(backend.invFilePath)
	inv = inv.Init()
	// ignore error for new file

	for _, name := range itemNames {
		item := inv.Take(name)

		for k, v := range labels {
			if v == "-" {
				delete(item, k)
			} else {
				item[k] = v
			}
		}

		if err := inv.Put(item); err != nil {
			return err
		}
	}

	return saveInventory(inv, backend.invFilePath)
}
