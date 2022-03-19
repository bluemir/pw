package backend

import (
	"context"
	"os"
	"os/exec"

	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/pw/pkg/util/console"
)

type RunOptions struct {
	Expr         string
	Labels       map[string]string
	ShortcutName string
	Template     string
	OutputFormat string
	WorkerNumber int
	Args         []string
}

func (backend *Backend) Run(opt *RunOptions) error {
	if opt.ShortcutName != "" {
		if v, ok := backend.inv.Shortcuts[opt.ShortcutName]; ok {
			opt.Expr = v
		}
	}
	if len(opt.Labels) > 0 && len(opt.Expr) > 0 {
		return errors.Errorf("`expr` and `label` option cannot use together")
	}

	items, err := backend.inv.ApplyExpr(opt.Expr)
	if err != nil {
		return err
	}
	logrus.Debug(items)

	// Output
	cout := console.New(os.Stdout)

	if opt.WorkerNumber <= 0 {
		opt.WorkerNumber = len(items)
	}

	t, ok := backend.inv.Templates[opt.Template]
	if !ok {
		return errors.Errorf("template not found(name: '%s'. check inventory file", opt.Template)
	}

	commandBuilder, err := NewCommandBuilder(t, opt.Args)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := workerpool.New(opt.WorkerNumber)
	// XXX support variant format
	formatter := newOutputFomatter(items, opt.OutputFormat)

	for _, item := range items {
		wp.Submit(executeCommand(
			ctx,
			item,
			cout,
			formatter,
			commandBuilder,
		))
	}
	wp.StopWait()

	cout.Write([]byte{'\n'})

	return nil
}
func getMaxItemNameLen(items []Item) int {
	max := 0
	for _, item := range items {
		if max < len(item["name"]) {
			max = len(item["name"])
		}
	}
	return max
}
func executeCommand(ctx context.Context, item Item, cout *console.Console, formatter OutputFormatter, commandBuilder *CommandBuilder) func() {
	return func() {
		defer cout.Close()

		cmd, err := commandBuilder.build(item)
		if err != nil {
			//TODO
			logrus.Fatal(err)
			return
		}

		logrus.Trace(cmd)

		c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)

		// TODO support variant format
		c.Stdout = cout.WithModifier(formatter.Modifier(item, "stdin"))
		c.Stderr = cout.WithModifier(formatter.Modifier(item, "stdout"))

		if err := c.Run(); err != nil {
			logrus.Fatal(err)
		}
	}
}
