package tui

type menuItem string

func (i menuItem) Title() string       { return string(i) }
func (i menuItem) Description() string { return string(i) }
func (i menuItem) FilterValue() string { return string(i) }
