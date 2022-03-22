package backend

import (
	"context"
	"os"

	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/pw/pkg/util/console"
	"github.com/bluemir/pw/pkg/util/exec"
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

	items, err := backend.inv.FindItem(opt.Expr)
	if err != nil {
		return err
	}
	logrus.Debug(items)

	// Output
	cout := console.New(os.Stdout)

	if opt.WorkerNumber <= 0 {
		opt.WorkerNumber = len(items)
	}

	t, ok := backend.inv.GetTemplate(opt.Template)
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
			logrus.Fatal(err)
			return
		}

		logrus.Trace(cmd)

		if err := exec.Exec(ctx, cmd,
			exec.WithStdout(cout.WithModifier(formatter.Modifier(item, "stdout"))),
			exec.WithStderr(cout.WithModifier(formatter.Modifier(item, "stderr"))),
		); err != nil {
			logrus.Fatal(err)
		}
	}
}
