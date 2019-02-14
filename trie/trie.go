package trie

type VisitFunc func([]byte, Value)
type Value int

type Trie interface {
	Set(k []byte, v Value)
}

type trie struct {
	root *node
}

func (t *trie) Each(f VisitFunc) {
	t.root.each(nil, f)
}

func (t *trie) set(key []byte, val Value) {
	if t.root == nil {
		t.root = makenode()
	}

	n := t.root

	for _, b := range key {
		if n.next[b] == nil {
			n.next[b] = makenode()
		}
		n = n.next[b]
	}

	n.v = val
	n.leaf = true
}
