package exec

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type ExecOption func(*exec.Cmd)

func WithWorkdir(path string) ExecOption {
	return func(cmd *exec.Cmd) {
		cmd.Dir = path
	}
}
func WithEnv(env []string) ExecOption {
	return func(cmd *exec.Cmd) {
		cmd.Env = append(cmd.Env, env...)
	}
}
func WithOut(out io.Writer) ExecOption {
	return func(cmd *exec.Cmd) {
		cmd.Stderr = out
		cmd.Stdout = out
	}
}
func WithStderr(out io.Writer) ExecOption {
	return func(cmd *exec.Cmd) {
		cmd.Stderr = out
	}

}
func WithStdout(out io.Writer) ExecOption {
	return func(cmd *exec.Cmd) {
		cmd.Stdout = out
	}
}

func Exec(ctx context.Context, commands []string, opt ...ExecOption) error {
	cmd := exec.CommandContext(ctx, commands[0], commands[1:]...)
	cmd.Env = os.Environ()

	for _, f := range opt {
		f(cmd)
	}

	return cmd.Run()
}
