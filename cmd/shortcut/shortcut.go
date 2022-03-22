package shortcut

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bluemir/pw/cmd/global"
	"github.com/bluemir/pw/pkg/backend"
)

func Register(cmd *kingpin.CmdClause, gOpts *global.Options) {
	{
		getCmd := cmd.Command("get", "get shortcut")
		conf := backend.GetShortcutOptions{}
		getCmd.Arg("name", "template name").
			StringsVar(&conf.Names)

		getCmd.Action(func(c *kingpin.ParseContext) error {
			backend, err := backend.New(gOpts.InventoryFile)
			if err != nil {
				return err
			}

			return backend.GetShortcut(&conf)
		})
	}
	{
		setCmd := cmd.Command("set", "set shortcut")
		conf := backend.SetShortcutOptions{}

		setCmd.Arg("name", "template name").
			Required().
			StringVar(&conf.Name)
		setCmd.Arg("args", "argument").
			Required().
			StringVar(&conf.Expr)
		setCmd.Action(func(c *kingpin.ParseContext) error {
			backend, _ := backend.New(gOpts.InventoryFile)

			return backend.SetShortcut(&conf)
		})

	}
	{
		delCmd := cmd.Command("delete", "delete shortcut").Alias("del").Alias("remove")
		conf := backend.DeleteShortcutOptions{Names: []string{}}

		delCmd.Arg("name", "template name").
			Required().
			StringsVar(&conf.Names)
		delCmd.Action(func(c *kingpin.ParseContext) error {

			backend, _ := backend.New(gOpts.InventoryFile)

			return backend.DeleteShortcut(&conf)
		})
	}
}
