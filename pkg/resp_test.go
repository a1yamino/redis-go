package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var testCases = "$6\r\nfoobar\r\n"

func TestRead(t *testing.T) {
	resp := NewResp(bufio.NewReader(strings.NewReader(testCases)))
	value, err := resp.Read()
	if err != nil {
		t.Fatal(err)
	}

	if value.typ != BULK {
		t.Errorf("expected %v, got %v", BULK, value.typ)
	}

	if value.str != "foobar" {
		t.Errorf("expected %v, got %v", "foobar", value.str)
	}
}

func TestWrite(t *testing.T) {
	resp := NewResp(bufio.NewReader(strings.NewReader(testCases)))
	value, err := resp.Read()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("value: %v\n", value)

	buf := &bytes.Buffer{}
	// writer := NewWriter(buf)

	t.Logf("buf: %v", buf.String())
}

func BenchmarkSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = []byte(fmt.Sprintf("%s %s", "foo", "bar"))
	}
}

func BenchmarkPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = []byte("foo" + " " + "bar")
	}
}

func BenchmarkWriteString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := &bytes.Buffer{}
		buf.WriteString("foo bar")
		_ = buf.Bytes()
	}
}

func BenchmarkJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = []byte(strings.Join([]string{"foo", "bar"}, " "))
	}
}
