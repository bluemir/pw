package console

import (
	"bufio"
	"bytes"
	"io"
)

// writer but has little tweak each line
type Console struct {
	output     io.Writer
	firstWrite *bool // TODO make atomic

	buf     *bytes.Buffer
	scanner *bufio.Scanner

	modifier []Modifier
}

type Modifier func(string) string

func New(output io.Writer) *Console {
	buf := bytes.NewBuffer(nil)

	scanner := bufio.NewScanner(buf)

	firstWrite := false

	return &Console{
		output:     output,
		firstWrite: &firstWrite,
		buf:        buf,
		scanner:    scanner,
		modifier:   nil,
	}
}

func (c *Console) Write(p []byte) (n int, err error) {
	n, err = c.buf.Write(p)
	if err != nil {
		return
	}

	for c.scanner.Scan() {
		if *c.firstWrite {
			_, err = c.output.Write([]byte{'\n'})
			if err != nil {
				return n, err
			}
		} else {
			*c.firstWrite = true
		}

		str := c.scanner.Text()
		for _, m := range c.modifier {
			str = m(str)
		}

		_, err = c.output.Write([]byte(str))
		if err != nil {
			return n, err
		}
	}

	return
}
func (c *Console) Close() error {
	// flush
	c.buf.WriteTo(c.output)
	return nil
}

func (c *Console) WithPrefix(prefix string) *Console {
	// share writer but print prefix each line.
	return c.WithModifier(func(in string) string {
		return prefix + in
	})
}
func (c *Console) WithModifier(m Modifier) *Console {
	// share writer but with modify each line.
	buf := bytes.NewBuffer(nil)

	scanner := bufio.NewScanner(buf)

	return &Console{
		output:     c.output,
		firstWrite: c.firstWrite,
		buf:        buf,
		scanner:    scanner,
		modifier:   append([]Modifier{m}, c.modifier...),
	}
}
