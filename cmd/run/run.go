package run

import (
	"github.com/bluemir/pw/cmd/global"
	"github.com/bluemir/pw/pkg/backend"
	"gopkg.in/alecthomas/kingpin.v2"
)

func Register(cmd *kingpin.CmdClause, opts *global.Options) {
	conf := backend.RunOptions{
		Labels: map[string]string{},
		Args:   []string{},
	}

	cmd.Flag("output", "display format(json, text, simple, wide or free format)").Short('o').
		Default("simple"). // wide, all, template
		StringVar(&conf.OutputFormat)
	cmd.Flag("expr", "condition that filter items").Short('e').
		StringVar(&conf.Expr)
	cmd.Flag("label", "label that filter items, shortcut of `expr`. cannot use together").
		StringMapVar(&conf.Labels)
	cmd.Flag("shotcut", "expr shortcut. if exist '--expr' will be ignored").Short('s').
		StringVar(&conf.ShortcutName)
	cmd.Flag("templates", "running template").Short('t').
		Default("default").
		StringVar(&conf.Template)
	cmd.Flag("worker", "number of worker. 0 is unlimited").Short('w').
		IntVar(&conf.WorkerNumber)
	cmd.Arg("args", "args to run").
		StringsVar(&conf.Args)
	cmd.Action(func(c *kingpin.ParseContext) error {
		backend, err := backend.New(opts.InventoryFile)
		if err != nil {
			return err
		}
		return backend.Run(&conf)
	})
}
