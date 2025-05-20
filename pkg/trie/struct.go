package trie

type Node struct {
	Char     rune
	Children map[rune]*Node
	IsEnd    bool
}

type Trie struct {
	Root *Node
}
