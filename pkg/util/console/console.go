package console

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

// writer but has little tweak each line
type Console struct {
	output  io.Writer
	written *bool
	mutex   *sync.Mutex

	buf     *bytes.Buffer
	scanner *bufio.Scanner

	modifier []Modifier
}

type Modifier func(string) string

func New(output io.Writer) *Console {
	buf := bytes.NewBuffer(nil)

	scanner := bufio.NewScanner(buf)

	written := false

	return &Console{
		output:   output,
		written:  &written,
		mutex:    &sync.Mutex{},
		buf:      buf,
		scanner:  scanner,
		modifier: nil,
	}
}

func (c *Console) Write(p []byte) (int, error) {
	if n, err := c.buf.Write(p); err != nil {
		return n, err
	}

	for c.scanner.Scan() {
		str := c.scanner.Text()
		for _, m := range c.modifier {
			str = m(str)
		}

		if _, err := c.writeWithMutex([]byte(str)); err != nil {
			return len(p), err
		}
	}
	return len(p), nil
}
func (c *Console) writeWithMutex(p []byte) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if *c.written {
		if _, err := c.output.Write([]byte{'\n'}); err != nil {
			return 0, err
		}
	} else {
		*c.written = true
	}

	return c.output.Write(p)
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
		output:   c.output,
		written:  c.written,
		mutex:    c.mutex,
		buf:      buf,
		scanner:  scanner,
		modifier: append([]Modifier{m}, c.modifier...),
	}
}
