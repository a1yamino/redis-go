package pkg

import (
	. "container/list"
	"testing"
)

func TestPushLeft(t *testing.T) {
	ql := &qlist{}
	ql.pushLeft("a")
	ql.pushLeft("b")
	ql.pushLeft("c")

	if ql.getLeft() != "c" {
		t.Errorf("expected c, got %s", ql.getLeft())
	}
	if ql.getRight() != "a" {
		t.Errorf("expected a, got %s", ql.getRight())
	}
}

func TestPushRight(t *testing.T) {
	ql := &qlist{}
	ql.pushRight("a")
	ql.pushRight("b")
	ql.pushRight("c")

	if ql.getLeft() != "a" {
		t.Errorf("expected a, got %s", ql.getLeft())
	}
	if ql.getRight() != "c" {
		t.Errorf("expected c, got %s", ql.getRight())
	}
}

func TestPush(t *testing.T) {
	ql := &qlist{}
	ql.pushLeft("a")
	ql.pushRight("b")
	ql.pushLeft("c")

	if ql.getLeft() != "c" {
		t.Errorf("expected c, got %s", ql.getLeft())
	}
	if ql.getRight() != "b" {
		t.Errorf("expected b, got %s", ql.getRight())
	}
}

func BenchmarkPush(b *testing.B) {
	ql := &qlist{}
	for i := 0; i < b.N; i++ {
		ql.pushLeft("a")
		ql.pushRight("b")
	}
}

func BenchmarkListPush(b *testing.B) {
	ll := New()
	for i := 0; i < b.N; i++ {
		ll.PushFront("a")
		ll.PushBack("b")
	}
}
