package template

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/bluemir/pw/cmd/global"
	"github.com/bluemir/pw/pkg/backend"
)

func Register(cmd *kingpin.CmdClause, gOpts *global.Options) {
	conf := backend.SetTemplateOptions{
		Name: "",
		Args: []string{},
	}

	//deleteFlag := false
	//cmd.Flag("delete", "delete template").BoolVar(&deleteFlag)

	cmd.Arg("name", "template name").
		StringVar(&conf.Name)
	cmd.Arg("args", "argument").
		StringsVar(&conf.Args)
	cmd.Action(func(c *kingpin.ParseContext) error {
		backend, _ := backend.New(gOpts.InventoryFile)
		// ignore error for new file
		return backend.SetTemplate(&conf)

	})

	/*
		{
			getCmd := cmd.Command("get", "get")

			getCmd.Action(func(c *kingpin.ParseContext) error {
				return nil
				//return client.Run(&conf)
			})
		}
		{
			setCmd := cmd.Command("set", "set item")
			conf := backend.TemplateSetOptions{
				Name: "",
				Args: []string{},
			}

			setCmd.Arg("name", "template name").
				StringVar(&conf.Name)
			setCmd.Arg("args", "argument").
				StringsVar(&conf.Args)
			setCmd.Action(func(c *kingpin.ParseContext) error {
				backend, _ := backend.New(gOpts.InventoryFile)
				// ignore error for new file
				return backend.SetTemplate(&conf)

			})
		}
		{
			delCmd := cmd.Command("del", "delelete item").Alias("delete").Alias("remove")
			conf := struct {
				TemplateNames []string
			}{[]string{}}

			delCmd.Arg("template", "templates").Required().
				StringsVar(&conf.TemplateNames)
			delCmd.Action(func(c *kingpin.ParseContext) error {
				return nil
			})
		}*/
}
