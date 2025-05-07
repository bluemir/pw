package console_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/bluemir/pw/pkg/util/console"
)

func TestSimpleConsole(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	console := console.New(buf)

	n, err := console.Write([]byte("aa"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, n)

	console.Close()

	assert.Equal(t, "aa\n", buf.String())
}
func TestMultiLineConsole(t *testing.T) {
	testCases := []string{
		"aa", "aa\naa", "bb\ncc",
	}

	for _, tc := range testCases {
		buf := bytes.NewBuffer(nil)

		console := console.New(buf)

		n, err := console.Write([]byte(tc))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(tc), n)

		console.Close()

		assert.Equal(t, tc+"\n", buf.String()) // always has last \n

	}
}
func TestWithPrefix(t *testing.T) {
	testCases := map[string]string{
		"aa":     "[hostname] aa\n",
		"aa\nbb": "[hostname] aa\n[hostname] bb\n",
	}
	for tc, expect := range testCases {
		buf := bytes.NewBuffer(nil)

		console := console.New(buf)
		console = console.WithPrefix("[hostname] ")

		n, err := console.Write([]byte(tc))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(tc), n)

		console.Close()

		assert.Equal(t, expect, buf.String())
	}
}
func TestWithMultiplePrefixedConsole(t *testing.T) {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)

	buf := bytes.NewBuffer(nil)

	console := console.New(buf)

	hostA := console.WithPrefix("[A] ")

	hostB := console.WithPrefix("[B] ")

	hostA.Write([]byte("aa"))
	hostB.Write([]byte("bb"))

	hostA.Close()
	hostB.Close()

	assert.Equal(t, "[A] aa\n[B] bb\n", buf.String())
}
func TestWithModifier(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	console := console.New(buf)

	hostA := console.WithModifier(func(in string) string {
		buf, _ := json.Marshal(map[string]string{
			"msg":  in,
			"host": "A",
		})

		return string(buf)
	})

	hostA.Write([]byte("aa"))

	hostA.Close()

	assert.Equal(t, `{"host":"A","msg":"aa"}`+"\n", buf.String())
}
func TestWithModifierChaining(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	console := console.New(buf).WithPrefix("[host]").WithPrefix("{server=api}")

	console.Write([]byte("aa"))

	console.Close()

	assert.Equal(t, `[host]{server=api}aa`+"\n", buf.String())
}
func TestWithModifierHasState(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	count := 0
	console := console.New(buf).WithModifier(func(str string) string {
		count++
		return fmt.Sprintf("%d | %s", count, str)
	})

	console.Write([]byte("aa\nbb"))

	console.Close()

	assert.Equal(t, "1 | aa\n2 | bb\n", buf.String())
}
