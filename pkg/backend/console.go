package backend

import "fmt"

type Console interface {
	Out(line uint, stream, text string, item Item) error
}

func NewConsole(format string, items []Item) Console {
	max := 0
	for _, item := range items {
		if max < len(item["name"]) {
			max = len(item["name"])
		}
	}
	switch format {
	default:
		return &TextConsole{}
	}
	return nil
}

type TextConsole struct {
}

func (f *TextConsole) Out(line uint, from, text string, item Item) error {
	fmt.Printf("%s\n", text)
	return nil
}
