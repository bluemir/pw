package backend

type SetTemplateOptions struct {
	Name string
	Args []string
}

func (backend *Backend) SetTemplate(opt *SetTemplateOptions) error {
	backend.inv.SetTemplate(opt.Name, opt.Args)

	return backend.Save()
}

type DeleteTemplateOptions struct {
	Names []string
}

func (backend *Backend) DeleteTemplate(opt *DeleteTemplateOptions) error {
	for _, name := range opt.Names {
		backend.inv.DeleteTemplate(name)
	}

	return backend.Save()
}
