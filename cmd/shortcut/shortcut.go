package shortcut

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bluemir/pw/cmd/global"
	"github.com/bluemir/pw/pkg/backend"
)

func Register(cmd *kingpin.CmdClause, gOpts *global.Options) {
	{
		getCmd := cmd.Command("get", "get")

		getCmd.Action(func(c *kingpin.ParseContext) error {
			return nil
		})
	}
	{
		setCmd := cmd.Command("set", "set item")
		conf := backend.SetShortcutOptions{}

		setCmd.Arg("name", "template name").
			StringVar(&conf.Name)
		setCmd.Arg("args", "argument").
			StringVar(&conf.Expr)
		setCmd.Action(func(c *kingpin.ParseContext) error {
			backend, _ := backend.New(gOpts.InventoryFile)

			return backend.SetShortcut(&conf)
		})

	}
	{
		delCmd := cmd.Command("del", "delelete item").Alias("delete").Alias("remove")
		conf := backend.DeleteShortcutOptions{Names: []string{}}

		delCmd.Arg("template", "templates").Required().
			StringsVar(&conf.Names)
		delCmd.Action(func(c *kingpin.ParseContext) error {

			backend, _ := backend.New(gOpts.InventoryFile)

			return backend.DeleteShortcut(&conf)
		})
	}
}
