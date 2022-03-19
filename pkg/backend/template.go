package backend

import (
	"strings"
)

type TemplateSetOptions struct {
	Name string
	Args []string
}

func (backend *Backend) SetTemplate(opt *TemplateSetOptions) error {

	backend.inv.Templates[opt.Name] = strings.Join(opt.Args, " ")

	return backend.Save()
}

type TmeplateDelOptions struct {
	Names []string
}

func (backend *Backend) TemplateDel(opt *TemplateSetOptions) error {
	delete(backend.inv.Templates, opt.Name)

	return backend.Save()
}
