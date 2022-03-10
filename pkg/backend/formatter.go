package backend

import (
	"encoding/json"
	"fmt"

	"github.com/bluemir/pw/pkg/util/console"
	"github.com/mgutz/str"
	"github.com/sirupsen/logrus"
)

type OutputFormatter interface {
	Modifier(item Item, from string) console.Modifier
}

func newOutputFomatter(items []Item, format string) OutputFormatter {
	max := 0
	for _, item := range items {
		if max < len(item["name"]) {
			max = len(item["name"])
		}
	}
	switch format {
	case "simple":
		return &SimpleFormatter{max}
	case "wide":
		return &WideFormatter{max}
	case "json":
		return &JsonFormatter{}
	case "text":
		return &TextFormatter{}
	default:
		logrus.Fatal("unknown format")
		return nil
	}
}

type TextFormatter struct {
}

func (c *TextFormatter) Modifier(item Item, from string) console.Modifier {
	return func(str string) string {
		return str
	}
}

type SimpleFormatter struct {
	max int
}

func (c *SimpleFormatter) Modifier(item Item, from string) console.Modifier {
	return func(text string) string {
		return fmt.Sprintf("%s | %s",
			str.PadLeft(item["name"], " ", c.max),
			text,
		)
	}
}

type WideFormatter struct {
	max int
}

func (c *WideFormatter) Modifier(item Item, from string) console.Modifier {
	line := 0
	return func(text string) string {
		line++

		return fmt.Sprintf("%s %05d %s | %s",
			str.PadLeft(item["name"], " ", c.max),
			line,
			from,
			text,
		)
	}
}

type JsonFormatter struct {
}

func (c *JsonFormatter) Modifier(item Item, stream string) console.Modifier {
	line := uint(0)
	return func(text string) string {
		line++

		buf, err := json.Marshal(struct {
			Line   uint
			Stream string
			Text   string
			Item   Item
		}{line, stream, text, item})

		if err != nil {
			logrus.Fatal(line)
		}

		return fmt.Sprintf("%s", buf)
	}
}
