package inventory

func (inv *Inventory) GetShortcut(name string) string {
	return inv.Shortcuts[name]
}
func (inv *Inventory) DeleteShortcut(name string) {
	delete(inv.Shortcuts, name)
}
func (inv *Inventory) SetShortcut(name string, expr string) {
	inv.Shortcuts[name] = expr
}
