package cache

type Node struct {
	next *Node
	prev *Node

	data interface{}
}

type List struct {
	sentinel *Node
	Len      int64
}

func NewList() *List {
	sentinel := &Node{}
	sentinel.prev = sentinel
	sentinel.next = sentinel

	return &List{sentinel: sentinel}
}

func (l *List) Insert(data interface{}) {
	n := &Node{data: data}

	n.next = l.sentinel.next
	l.sentinel.next.prev = n

	l.sentinel.next = n
	n.prev = l.sentinel

	l.Len += 1
}

func (l *List) search(data interface{}) (*Node, bool) {
	x := l.sentinel.next
	for {
		if x != l.sentinel && x.data != data {
			x = x.next
			continue
		}

		if x == l.sentinel {
			return nil, false
		}
		return x, true
	}
}

func (l *List) Contains(data interface{}) bool {
	if _, ok := l.search(data); !ok {
		return false
	}

	return true
}

func (l *List) Delete(data interface{}) {
	n, ok := l.search(data)
	if !ok {
		return
	}

	l.delete(n)
}

func (l *List) delete(n *Node) {
	n.prev.next = n.next
	n.next.prev = n.prev

	l.Len -= 1
}

func (l *List) DeleteHead() (data interface{}) {
	data = l.sentinel.next.data
	l.delete(l.sentinel.next)

	return data
}

func (l *List) DeleteTail() (data interface{}) {
	data = l.sentinel.prev.data
	l.delete(l.sentinel.prev)

	return data
}

func (l *List) Range() []interface{} {
	x := l.sentinel.next

	var nodes []interface{}
	for {
		if x == l.sentinel {
			break
		}
		nodes = append(nodes, x.data)
		x = x.next
	}

	return nodes
}
