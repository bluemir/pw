package backend

type SetOptions struct {
	ItemNames []string
	Labels    map[string]string
}

func (backend *Backend) Set(opt *SetOptions) error {
	inv, _ := loadInventory(backend.invFilePath)
	inv = inv.Init()
	// ignore error for new file

	for _, name := range opt.ItemNames {
		item := inv.Take(name)

		for k, v := range opt.Labels {
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
