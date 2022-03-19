package del

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bluemir/pw/cmd/global"
	"github.com/bluemir/pw/pkg/backend"
)

func Register(cmd *kingpin.CmdClause, gOpts *global.Options) {
	cmd = cmd.Alias("del").Alias("remove")
	conf := backend.DeleteItemOptions{
		ItemNames: []string{},
	}

	cmd.Arg("items", "items").Required().
		StringsVar(&conf.ItemNames)
	cmd.Action(func(c *kingpin.ParseContext) error {
		backend, err := backend.New(gOpts.InventoryFile)
		if err != nil {
			return err
		}
		return backend.DeleteItems(&conf)
	})
}
