package pkg

type qlist struct {
	head *qnode
	tail *qnode

	len int
}

type qnode struct {
	prev *qnode
	next *qnode

	len  int
	data []string

	direction bool // true for left, false for right
}

const (
	Left   = true
	Right  = false
	maxLen = 512
)

func (ql *qlist) pushLeft(data string) {
	h := ql.head
	if h == nil {
		ql.head = &qnode{data: []string{data}, direction: Left, len: 1}
		ql.len++
		ql.tail = ql.head
		return
	}

	if h.direction == Right || h.len >= maxLen {
		n := &qnode{data: []string{data}, direction: Left, len: 1}
		n.next = h
		h.prev = nil
		ql.len++
		ql.head = n
	}

	if h.direction == Left {
		h.data = append(h.data, data)
		h.len++
		ql.len++
		return
	}
}

func (ql *qlist) pushRight(data string) {
	t := ql.tail
	if t == nil {
		ql.tail = &qnode{data: []string{data}, direction: Right, len: 1}
		ql.head = ql.tail
		ql.len++
		return
	}

	if t.direction == Left || t.len >= maxLen {
		n := &qnode{data: []string{data}, direction: Right, len: 1}
		n.prev = t
		t.next = nil
		ql.tail = n
		ql.len++
	}

	if t.direction == Right {
		t.data = append(t.data, data)
		t.len++
		ql.len++
		return
	}
}

func (ql *qlist) getLeft() string {
	h := ql.head
	if h == nil {
		return ""
	}

	if h.len == 0 {
		return ""
	}

	var v string

	if h.direction == Right {
		v = h.data[0]
	} else {
		v = h.data[h.len-1]
	}
	return v
}

func (ql *qlist) getRight() string {
	t := ql.tail
	if t == nil {
		return ""
	}

	if t.len == 0 {
		return ""
	}

	var v string

	if t.direction == Left {
		v = t.data[0]
	} else {
		v = t.data[t.len-1]
	}
	return v
}

func (ql *qlist) getRange(start, stop int) []string {
	h := ql.head
	if h == nil {
		return nil
	}

	for h.len <= start {
		h = h.next
		start -= h.len
		stop -= h.len
	}

	res := make([]string, 0, stop-start+1)

	for i := start; i <= stop; i++ {
		res = append(res, h.data[i])
		if i == h.len-1 {
			h = h.next
			i = 0
			stop -= h.len
		}
	}
	return res
}

func (ql *qlist) popLeft() string {
	h := ql.head
	if h == nil {
		return ""
	}

	if h.len == 0 {
		return ""
	}

	var v string

	if h.direction == Right {
		v = h.data[0]
		h.data = h.data[1:]
		h.len--
		ql.len--
		if h.len == 0 {
			ql.head = h.next
		}
	} else {
		v = h.data[h.len-1]
		h.data = h.data[:h.len-1]
		h.len--
		ql.len--
		if h.len == 0 {
			ql.head = h.next
		}
	}

	return v
}

func (ql *qlist) popRight() string {
	t := ql.tail
	if t == nil {
		return ""
	}

	if t.len == 0 {
		return ""
	}

	var v string

	if t.direction == Left {
		v = t.data[0]
		t.data = t.data[1:]
		t.len--
		ql.len--
		if t.len == 0 {
			ql.tail = t.prev
		}
	} else {
		v = t.data[t.len-1]
		t.data = t.data[:t.len-1]
		t.len--
		ql.len--
		if t.len == 0 {
			ql.tail = t.prev
		}
	}

	return v
}
