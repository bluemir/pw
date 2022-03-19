package get

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bluemir/pw/cmd/global"
	"github.com/bluemir/pw/pkg/backend"
)

func Register(cmd *kingpin.CmdClause, gOpts *global.Options) {
	conf := backend.GetOptions{}

	cmd.Flag("output", "item display format").Short('o').
		Default("name"). // name, yaml, json, wide
		StringVar(&conf.OutputFormat)
	cmd.Flag("expr", "item filter").Short('e').
		StringVar(&conf.Expr)
	cmd.Flag("shotcut", "expr shortcut. if exist '--expr' will be ignored").Short('s').
		StringVar(&conf.ShortcutName)
	cmd.Action(func(c *kingpin.ParseContext) error {
		backend, err := backend.New(gOpts.InventoryFile)
		if err != nil {
			return err
		}
		return backend.GetItems(&conf)
	})
}
