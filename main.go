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
	LogLevel         int
	InventoryFile    string
	OutputFormat     string
	Expr             string
	ShortcutExprName string
	Template         string
	Labels           map[string]string
	Args             []string
	ItemNames        []string
	WorkerNumber     int
}

// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md

func main() {

	conf := &Config{Labels: map[string]string{}}

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
		StringVar(&conf.OutputFormat)
	run.Flag("expr", "condition that filter items").Short('e').
		StringVar(&conf.Expr)
	run.Flag("shotcut", "expr shortcut. if exist '--expr' will be ignored").Short('s').
		StringVar(&conf.ShortcutExprName)
	run.Flag("templates", "running template").Short('t').
		Default("default").
		StringVar(&conf.Template)
	run.Flag("worker", "number of worker. 0 is unlimited").Short('w').
		IntVar(&conf.WorkerNumber)
	run.Arg("args", "args to run").
		StringsVar(&conf.Args)

	set := cli.Command("set", "Set item")
	set.Flag("label", "labels").Short('l').
		StringMapVar(&conf.Labels)
	set.Arg("item", "items").Required().
		StringsVar(&conf.ItemNames)

	get := cli.Command("get", "Get item")
	get.Flag("output", "item display format").Short('o').
		Default("name"). // name, yaml, json, wide
		StringVar(&conf.OutputFormat)
	get.Flag("expr", "item filter").Short('e').
		StringVar(&conf.Expr)
	get.Flag("shotcut", "expr shortcut. if exist '--expr' will be ignored").Short('s').
		StringVar(&conf.ShortcutExprName)

	del := cli.Command("del", "Delete item")
	del.Arg("item", "items").Required().
		StringsVar(&conf.ItemNames)

	template := cli.Command("template", "Set template")
	template.Arg("name", "template name").
		StringVar(&conf.Template)
	template.Arg("args", "argument").
		StringsVar(&conf.Args)

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
		if err := b.Run(conf.Expr, conf.ShortcutExprName, conf.Template, conf.OutputFormat, conf.WorkerNumber, conf.Args); err != nil {
			logrus.Fatal(err)
		}
	case set.FullCommand():
		if err := b.Set(conf.ItemNames, conf.Labels); err != nil {
			logrus.Fatal(err)
		}
	case get.FullCommand():
		if err := b.Get(conf.Expr, conf.ShortcutExprName, conf.OutputFormat); err != nil {
			logrus.Fatal(err)
		}
	case del.FullCommand():
		if err := b.Del(conf.ItemNames); err != nil {
			logrus.Fatal(err)
		}
	case template.FullCommand():
		if err := b.Template(conf.Template, conf.Args); err != nil {
			logrus.Fatal(err)
		}
	}

}
