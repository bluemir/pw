package console_test

import (
	"bytes"
	"encoding/json"
	"testing"

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

	assert.Equal(t, "aa", buf.String())
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

		assert.Equal(t, tc, buf.String())

	}
}
func TestWithPrefix(t *testing.T) {
	testCases := map[string]string{
		"aa":     "[hostname] aa",
		"aa\nbb": "[hostname] aa\n[hostname] bb",
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
	buf := bytes.NewBuffer(nil)

	console := console.New(buf)

	hostA := console.WithPrefix("[A] ")

	hostB := console.WithPrefix("[B] ")

	hostA.Write([]byte("aa"))
	hostB.Write([]byte("bb"))

	console.Close()

	assert.Equal(t, "[A] aa\n[B] bb", buf.String())
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

	console.Close()

	assert.Equal(t, `{"host":"A","msg":"aa"}`, buf.String())
}
func TestWithModifierChaining(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	console := console.New(buf).WithPrefix("[host]").WithPrefix("{server=api}")

	console.Write([]byte("aa"))

	console.Close()

	assert.Equal(t, `[host]{server=api}aa`, buf.String())
}
