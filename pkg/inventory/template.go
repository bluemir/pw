package inventory

func (inv *Inventory) GetTemplate(name string) ([]string, bool) {
	c, ok := inv.Templates[name]
	return c, ok
}
func (inv *Inventory) SetTemplate(name string, args []string) {
	inv.Templates[name] = args
}
func (inv *Inventory) DeleteTemplate(name string) {
	delete(inv.Templates, name)
}
