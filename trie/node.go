package trie

type Node interface {
}

type node struct {
	next [256]*node
	v    Value
	leaf bool
}

func makenode() *node {
	return &node{}
}

func (n *node) each(prefix []byte, f VisitFunc) {
	if n == nil {
		return
	}

	// we assume a node has an implicit leaf
	if n.leaf {
		f(prefix, n.v)
	}

	prefix = append(prefix, 0)
	for b, nn := range n.next {
		prefix[len(prefix)-1] = byte(b)
		nn.each(prefix, f)
	}
}

func (n *node) set(key []byte, val Value) {
	if n == nil {
		n = makenode()
	}

	for _, b := range key {
		if n.next[b] == nil {
			n.next[b] = makenode()
		}
		n = n.next[b]
	}

	n.v = val
	n.leaf = true
}

func (n *node) get(key []byte) (*Value, bool) {
	if n == nil {
		return nil, false
	}

	for _, b := range key {
		if n.next[b] == nil {
			return nil, false
		}

		n = n.next[b]
	}

	return &n.v, true
}
