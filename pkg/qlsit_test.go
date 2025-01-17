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

func TestPopLeft(t *testing.T) {
	ql := &qlist{}
	ql.pushLeft("a")
	ql.pushLeft("b")
	ql.pushLeft("c")

	if s := ql.popLeft(); s != "c" {
		t.Errorf("expected c, got %s", s)
	}
	if s := ql.popLeft(); s != "b" {
		t.Errorf("expected b, got %s", s)
	}
	if s := ql.popLeft(); s != "a" {
		t.Errorf("expected a, got %s", s)
	}
	if s := ql.popLeft(); s != "" {
		t.Errorf("expected empty, got %s", s)
	}
}

func TestPopRight(t *testing.T) {
	ql := &qlist{}
	ql.pushRight("a")
	ql.pushRight("b")
	ql.pushRight("c")

	if s := ql.popRight(); s != "c" {
		t.Errorf("expected c, got %s", s)
	}
	if s := ql.popRight(); s != "b" {
		t.Errorf("expected b, got %s", s)
	}
	if s := ql.popRight(); s != "a" {
		t.Errorf("expected a, got %s", s)
	}
	if s := ql.popRight(); s != "" {
		t.Errorf("expected empty, got %s", s)
	}
}

func TestPop(t *testing.T) {
	ql := &qlist{}
	ql.pushLeft("a")
	ql.pushRight("b")
	ql.pushLeft("c")

	if s := ql.popLeft(); s != "c" {
		t.Errorf("expected c, got %s", s)
	}
	if s := ql.popRight(); s != "b" {
		t.Errorf("expected b, got %s", s)
	}
	if s := ql.popLeft(); s != "a" {
		t.Errorf("expected a, got %s", s)
	}
	if s := ql.popRight(); s != "" {
		t.Errorf("expected empty, got %s", s)
	}
}

func TestGetRange(t *testing.T) {
	ql := &qlist{}
	ql.pushLeft("a")
	ql.pushLeft("b")
	ql.pushLeft("c")
	// ql.pushRight("d")
	// ql.pushRight("e")
	// c a b d
	s := ql.getRange(0, 2)

	for i, v := range s {
		t.Logf("%d: %v", i, v)
	}
}

func (qr qrange) String() string {
	res := ""
	if qr.direction == Left {
		for i := len(qr.data) - 1; i >= 0; i-- {
			res += qr.data[i]
		}
	} else {
		for i := 0; i < len(qr.data); i++ {
			res += qr.data[i]
		}
	}
	return res
}
