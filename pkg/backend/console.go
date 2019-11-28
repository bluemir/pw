package backend

import (
	"encoding/json"
	"fmt"

	"github.com/mgutz/str"
	"github.com/sirupsen/logrus"
)

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
	case "simple":
		return &SimpleConsole{max}
	case "wide":
		return &WideConsole{max}
	case "json":
		return &JsonConsole{}
	case "text":
		return &TextConsole{}
	default:
		logrus.Fatal("unknown format")
		return nil
	}
	return nil
}

type TextConsole struct {
}

func (c *TextConsole) Out(line uint, from, text string, item Item) error {
	fmt.Printf("%s\n", text)
	return nil
}

type SimpleConsole struct {
	max int
}

func (c *SimpleConsole) Out(line uint, from, text string, item Item) error {
	fmt.Printf("%s | %s\n",
		str.PadLeft(item["name"], " ", c.max),
		text,
	)
	return nil
}

type WideConsole struct {
	max int
}

func (c *WideConsole) Out(line uint, from, text string, item Item) error {
	fmt.Printf("%s %05d %s | %s\n",
		str.PadLeft(item["name"], " ", c.max),
		line,
		from,
		text,
	)
	return nil
}

type JsonConsole struct {
}

func (c *JsonConsole) Out(line uint, stream, text string, item Item) error {
	buf, err := json.Marshal(struct {
		Line   uint
		Stream string
		Text   string
		Item   Item
	}{line, stream, text, item})
	fmt.Printf("%s\n", buf)
	return err
}
