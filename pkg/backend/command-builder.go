package backend

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

type CommandBuilder struct {
	parts []*template.Template
	args  []string
}

func NewCommandBuilder(tmplCommand []string, args []string) (*CommandBuilder, error) {
	logrus.Tracef("args: %+v", args)

	parts := []string{}
	for _, part := range tmplCommand {
		if strings.EqualFold(part, "{{args}}") {
			parts = append(parts, args...)
		} else {
			parts = append(parts, part)
		}
	}
	logrus.Tracef("template parts: %+v", parts)

	fMap := map[string]interface{}{
		"args": func() string {
			return strings.Join(args, " ")
		},
		"arg": func(index int) string {
			logrus.Tracef("arg index: %d", index)
			if index < 1 || index > len(args) {
				logrus.Warnf("index of bound: {{arg %d}}", index)
				return ""
			}
			return args[index-1]
		},
	}

	tParts := []*template.Template{}
	for n, part := range parts {
		logrus.Tracef("part: %s", part)
		tmpl, err := template.New("__").Funcs(fMap).Parse(part)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		tParts = append(tParts, tmpl)
		logrus.Debugf("%d %s", n, part)
	}

	return &CommandBuilder{tParts, args}, nil
}

func (cb *CommandBuilder) build(item Item) ([]string, error) {
	parts := []string{}

	for _, tmpl := range cb.parts {
		var cmd bytes.Buffer
		err := tmpl.Execute(&cmd, item)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		parts = append(parts, cmd.String())
	}
	return parts, nil
}
