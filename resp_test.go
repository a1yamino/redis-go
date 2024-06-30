package main

import (
	"strings"
	"testing"
)

var testCases = "$6\r\nfoobar\r\n"

func TestRead(t *testing.T) {
	resp := NewResp(strings.NewReader(testCases))
	value, err := resp.Read()
	if err != nil {
		t.Fatal(err)
	}

	if value.typ != "bulk" {
		t.Errorf("expected %v, got %v", "bulk", value.typ)
	}

	if value.bulk != "foobar" {
		t.Errorf("expected %v, got %v", "foobar", value.bulk)
	}
}
