package console

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// writer but has little tweak each line
type Console struct {
	output io.Writer

	buf *bytes.Buffer

	modifier []Modifier
}

type Modifier func(string) string

func New(output io.Writer) *Console {
	buf := bytes.NewBuffer(nil)

	return &Console{
		output:   output,
		buf:      buf,
		modifier: nil,
	}
}

func (c *Console) Write(p []byte) (int, error) {
	if n, err := c.buf.Write(p); err != nil {
		return n, err
	}

	return len(p), c.handle(false)
}

func (c *Console) Close() error {
	return c.handle(true)
}
func (c *Console) Flush() {
	if err := c.handle(true); err != nil {
		logrus.Error(err)
	}
}

func (c *Console) handle(flush bool) error {
	for {
		line, err := c.buf.ReadString('\n')
		logrus.Trace(line, err, flush)
		if err != nil {
			if err == io.EOF && !flush {
				//c.buf.Reset()
				if _, err := c.buf.WriteString(line); err != nil {
					return errors.WithStack(err)
				}

				return nil
			}
		}
		for _, m := range c.modifier {
			line = m(line)
		}

		if _, err := c.output.Write([]byte(line)); err != nil {
			return err
		}
		if flush {
			c.output.Write([]byte{'\n'})
			return nil
		}
	}
}
func (c *Console) WithPrefix(prefix string) *Console {
	// share writer but print prefix each line.
	return c.WithModifier(func(in string) string {
		return prefix + in
	})
}
func (c *Console) WithSuffix(subfix string) *Console {
	return c.WithModifier(func(in string) string {
		return in + subfix
	})
}
func (c *Console) WithModifier(m Modifier) *Console {
	// share writer but with modify each line.
	buf := bytes.NewBuffer(nil)

	return &Console{
		output:   c.output,
		buf:      buf,
		modifier: append([]Modifier{m}, c.modifier...),
	}
}
