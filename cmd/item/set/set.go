package set

import (
	"github.com/alecthomas/kingpin/v2"

	"github.com/bluemir/pw/cmd/global"
	"github.com/bluemir/pw/pkg/backend"
)

func Register(cmd *kingpin.CmdClause, gOpts *global.Options) {
	conf := backend.SetItemOptions{
		Labels:    map[string]string{},
		ItemNames: []string{},
	}

	cmd.Flag("label", "labels").Short('l').
		StringMapVar(&conf.Labels)
	cmd.Arg("item", "items").Required().
		StringsVar(&conf.ItemNames)

	cmd.Action(func(c *kingpin.ParseContext) error {
		backend, err := backend.New(gOpts.InventoryFile)
		if err != nil {
			return err
		}
		return backend.SetItems(&conf)
	})
}
