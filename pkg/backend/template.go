package backend

type TemplateSetOptions struct {
	Name string
	Args []string
}

func (backend *Backend) SetTemplate(opt *TemplateSetOptions) error {
	backend.inv.SetTemplate(opt.Name, opt.Args)

	return backend.Save()
}

type DeleteTemplateOptions struct {
	Names []string
}

func (backend *Backend) TemplateDel(opt *DeleteTemplateOptions) error {
	for _, name := range opt.Names {
		backend.inv.DeleteTemplate(name)
	}

	return backend.Save()
}
