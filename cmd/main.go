package cmd

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/pw/cmd/global"
	delCmd "github.com/bluemir/pw/cmd/item/delete"
	getCmd "github.com/bluemir/pw/cmd/item/get"
	setCmd "github.com/bluemir/pw/cmd/item/set"
	runCmd "github.com/bluemir/pw/cmd/run"
	shortcutCmd "github.com/bluemir/pw/cmd/shortcut"
	templateCmd "github.com/bluemir/pw/cmd/template"
	"github.com/bluemir/pw/internal/buildinfo"
)

const (
	describe        = `run parallel worker with template`
	defaultLogLevel = logrus.WarnLevel
)

func Run() error {
	conf := struct {
		logLevel  int
		logFormat string
	}{}

	globalOpts := global.Options{}

	app := kingpin.New(buildinfo.AppName, describe)
	app.Version(buildinfo.Version + "/" + buildinfo.BuildTime)

	app.Flag("verbose", "Log level").
		Short('v').
		CounterVar(&conf.logLevel)
	app.Flag("log-format", "Log format").
		PlaceHolder("FORMAT").
		StringVar(&conf.logFormat)
	app.PreAction(func(*kingpin.ParseContext) error {
		level := logrus.Level(conf.logLevel) + defaultLogLevel
		logrus.SetOutput(os.Stderr)
		logrus.SetLevel(level)
		logrus.SetReportCaller(true)
		logrus.Infof("logrus level: %s", level)

		switch conf.logFormat {
		case "text-color":
			logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		case "text":
			logrus.SetFormatter(&logrus.TextFormatter{})
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		case "":
			// do nothing. it means smart.
		default:
			return errors.Errorf("unknown log format")
		}

		return nil
	})

	app.Flag("inventory", "inventory file. contain items and templates").Short('i').
		PlaceHolder("$HOME/.inventory.yaml").Default(os.ExpandEnv("$HOME/.inventory.yaml")).
		Envar("INVENTORY_FILE").
		StringVar(&globalOpts.InventoryFile)

	runCmd.Register(app.Command("run", "run command"), &globalOpts)
	getCmd.Register(app.Command("get", "get item"), &globalOpts)
	setCmd.Register(app.Command("set", "set item"), &globalOpts)
	delCmd.Register(app.Command("delete", "delete item"), &globalOpts)
	templateCmd.Register(app.Command("template", "set template"), &globalOpts)
	shortcutCmd.Register(app.Command("shortcut", "shortcut commands"), &globalOpts)

	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	logrus.Debug(cmd)
	return nil
}
