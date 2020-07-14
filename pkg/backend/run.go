package backend

import (
	"bufio"
	"context"
	"io"
	"os/exec"

	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func (backend *Backend) Run(expr string, shortcut string, template string, outputFormat string, workerNum int, args []string) error {
	inv, err := loadInventory(backend.invFilePath)
	if err != nil {
		return err
	}
	inv = inv.Init()

	if shortcut != "" {
		if v, ok := inv.Init().Shortcuts[shortcut]; ok {
			expr = v
		}
	}

	items, err := inv.ApplyExpr(expr)
	if err != nil {
		return err
	}
	logrus.Debug(items)

	// Output

	console := NewConsole(outputFormat, items)

	if workerNum <= 0 {
		workerNum = len(items)
	}

	t, ok := inv.Templates[template]
	if !ok {
		return errors.Errorf("template not found(name: '%s'. check inventory file", template)
	}

	commandBuilder, err := NewCommandBuilder(t, args)
	if err != nil {
		return err
	}

	wp := workerpool.New(workerNum)
	for _, item := range items {
		wp.Submit(executeCommand(context.TODO(), item, console, commandBuilder))
	}
	wp.StopWait()

	return nil
}
func executeCommand(ctx context.Context, item Item, console Console, commandBuilder *CommandBuilder) func() {
	return func() {
		cmd, err := commandBuilder.build(item)
		if err != nil {
			//TODO
			logrus.Fatal(err)
			return
		}

		logrus.Trace(cmd)

		c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
		defer c.Wait()

		stdout, err := c.StdoutPipe()
		if err != nil {
			logrus.Warn("pipe stdout", err)
			logrus.Fatal(err)
			return
		}
		stderr, err := c.StderrPipe()
		if err != nil {
			logrus.Warn("pipe stderr", err)
			logrus.Fatal(err)
			return
		}
		eg, _ := errgroup.WithContext(context.Background())
		eg.Go(read(item, "stdout", stdout, console))
		eg.Go(read(item, "stderr", stderr, console))

		eg.Go(c.Start)
		eg.Wait()
	}
}

func read(item Item, from string, reader io.Reader, console Console) func() error {
	name := item["_"]
	return func() error {
		ln := uint(1)
		r := bufio.NewScanner(reader)
		for r.Scan() {
			logrus.Debugf("read line from %s %s", name, from)
			console.Out(ln, from, r.Text(), item)
			ln++
		}
		if err := r.Err(); err != nil {
			logrus.Errorf("read line error %s %s: %q", name, from, err)
			return err
		}
		logrus.Debugf("end of stream")
		return nil
	}
}
