package main

import (
	"os"

	"github.com/codingconcepts/env"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bluemir/pw/pkg/backend"
)

var VERSION string

type Config struct {
	LogLevel      int
	InventoryFile string

	RunOptions         backend.RunOptions
	GetOptions         backend.GetOptions
	SetOptions         backend.SetOptions
	DelOptions         backend.DelOptions
	TemplateSetOptions backend.TemplateSetOptions
}

// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md

func main() {
	conf := &Config{
		SetOptions: backend.SetOptions{
			Labels: map[string]string{},
		},
	}

	if err := env.Set(conf); err != nil {
		logrus.Fatal(err)
		return
	}

	cli := kingpin.New("pw", "run parallel worker with template")

	cli.Flag("verbose", "log level").Short('v').
		CounterVar(&conf.LogLevel)

	cli.Flag("inventory", "inventory file. contain items and templates").Short('i').
		PlaceHolder("$HOME/.inventory.yaml").Default(os.ExpandEnv("$HOME/.inventory.yaml")).
		StringVar(&conf.InventoryFile)

	run := cli.Command("run", "Run command")
	run.Flag("output", "display format(json, text, simple, wide or free format)").Short('o').
		Default("simple"). // wide, all, template
		StringVar(&conf.RunOptions.OutputFormat)
	run.Flag("expr", "condition that filter items").Short('e').
		StringVar(&conf.RunOptions.Expr)
	run.Flag("shotcut", "expr shortcut. if exist '--expr' will be ignored").Short('s').
		StringVar(&conf.RunOptions.ShortcutExprName)
	run.Flag("templates", "running template").Short('t').
		Default("default").
		StringVar(&conf.RunOptions.Template)
	run.Flag("worker", "number of worker. 0 is unlimited").Short('w').
		IntVar(&conf.RunOptions.WorkerNumber)
	run.Arg("args", "args to run").
		StringsVar(&conf.RunOptions.Args)

	set := cli.Command("set", "Set item")
	set.Flag("label", "labels").Short('l').
		StringMapVar(&conf.SetOptions.Labels)
	set.Arg("item", "items").Required().
		StringsVar(&conf.SetOptions.ItemNames)

	get := cli.Command("get", "Get item")
	get.Flag("output", "item display format").Short('o').
		Default("name"). // name, yaml, json, wide
		StringVar(&conf.GetOptions.OutputFormat)
	get.Flag("expr", "item filter").Short('e').
		StringVar(&conf.GetOptions.Expr)
	get.Flag("shotcut", "expr shortcut. if exist '--expr' will be ignored").Short('s').
		StringVar(&conf.GetOptions.ShortcutExprName)

	del := cli.Command("del", "Delete item")
	del.Arg("item", "items").Required().
		StringsVar(&conf.DelOptions.ItemNames)

	template := cli.Command("template", "Set template")

	template.Arg("name", "template name").
		StringVar(&conf.TemplateSetOptions.Name)
	template.Arg("args", "argument").
		StringsVar(&conf.TemplateSetOptions.Args)

	cli.Version(VERSION)

	cmd := kingpin.MustParse(cli.Parse(os.Args[1:]))

	// adjust loglevel
	level := logrus.Level(conf.LogLevel) + logrus.ErrorLevel
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(level) // error level is default
	logrus.Infof("error level: %s", level)

	b, err := backend.New(conf.InventoryFile)
	if err != nil {
		logrus.Fatal(err)
	}

	switch cmd {
	case run.FullCommand():
		if err := b.Run(&conf.RunOptions); err != nil {
			logrus.Fatal(err)
		}
	case set.FullCommand():
		if err := b.Set(&conf.SetOptions); err != nil {
			logrus.Fatal(err)
		}
	case get.FullCommand():
		if err := b.Get(&conf.GetOptions); err != nil {
			logrus.Fatal(err)
		}
	case del.FullCommand():
		if err := b.Del(&conf.DelOptions); err != nil {
			logrus.Fatal(err)
		}
	case template.FullCommand():
		if err := b.TemplateSet(&conf.TemplateSetOptions); err != nil {
			logrus.Fatal(err)
		}

		// TODO
		// pw template set
		// pw template del
		// pw shortcut set
		// pw shortcut del
	}

}
