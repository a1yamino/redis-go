package pkg

import "math/rand"

type skipList struct {
	header *skipListNode
	tail   *skipListNode
	length uint64
	level  int
}

func newSkipList() *skipList {
	return &skipList{
		header: &skipListNode{
			next: make([]*skipListLevel, 32),
		},
		level: 1,
	}

}

type skipListNode struct {
	next  []*skipListLevel
	bwd   *skipListNode
	str   string
	score float64
}

type skipListLevel struct {
	forward *skipListNode
	span    int
}

// compareSkipListNodes compares two skip list nodes.
// return true if a < b, otherwise return false.
func compareSkipListNodes(a, b *skipListNode) bool {
	if a.score != b.score {
		return a.score < b.score
	}
	return a.str < b.str
}

// generate a random level for skip list via coin flip.
func randomLevel() int {
	level := 1
	for (rand.Int31()&0xFFFF)%2 == 0 {
		level++
	}
	return level
}

func (sl *skipList) Insert(score float64, str string) *skipListNode {
	// create a new node
	node := &skipListNode{
		str:   str,
		score: score,
	}

	// find the insert position
	update := make([]*skipListNode, 32)
	rank := make([]int, 32)
	x := sl.header
	for i := sl.level - 1; i >= 1; i-- {
		rank[i] = rank[i-1]

		for x.next[i-1] != nil && compareSkipListNodes(x.next[i-1].forward, node) {
			rank[i] += x.next[i-1].span
			x = x.next[i-1].forward
		}
		update[i] = x
	}

	// determine the level of the new node
	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.header
			rank[i] = 0
		}
		sl.level = level
	}

	// insert the new node
	node.next = make([]*skipListLevel, level)
	for i := 1; i <= level; i++ {
		node.next[i-1].forward = update[i].next[i-1].forward
		update[i].next[i-1].forward = node
		if i == 1 {
			node.bwd = update[i]
		}
		if node.next[i-1].forward != nil {
			node.next[i-1].forward.bwd = node
		}
		node.next[i-1].span = update[i].next[i-1].span - (rank[1] - rank[i])
		update[i].next[i-1].span = (rank[1] - rank[i]) + 1
	}

	sl.length++
	return node
}

func (sl *skipList) Delete(score float64, str string) bool {
	// find the node
	update := make([]*skipListNode, 32)
	x := sl.header
	for i := sl.level - 1; i >= 1; i-- {
		for x.next[i-1] != nil && compareSkipListNodes(x.next[i-1].forward, &skipListNode{str: str, score: score}) {
			x = x.next[i-1].forward
		}
		update[i] = x
	}

	x = x.next[0].forward
	if x == nil || x.score != score || x.str != str {
		return false
	}

	// delete the node
	for i := 1; i <= sl.level; i++ {
		if update[i].next[i-1].forward == x {
			update[i].next[i-1].forward = x.next[i-1].forward
			update[i].next[i-1].span += x.next[i-1].span - 1
		} else {
			update[i].next[i-1].span--
		}
	}

	if x.next[0].forward != nil {
		x.next[0].forward.bwd = x.bwd
	}

	for sl.level > 1 && sl.header.next[sl.level-1] == nil {
		sl.level--
	}

	sl.length--
	return true
}
