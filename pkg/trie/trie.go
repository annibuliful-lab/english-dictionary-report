package trie

func NewNode(char rune) *Node {
	return &Node{
		Char:     char,
		Children: make(map[rune]*Node),
		IsEnd:    false,
	}
}

func NewTrie() *Trie {
	return &Trie{
		Root: NewNode(0),
	}
}

func (t *Trie) Insert(word string) {
	current := t.Root

	for _, char := range word {
		if _, exists := current.Children[char]; !exists {
			current.Children[char] = NewNode(char)
		}

		current = current.Children[char]
	}

	current.IsEnd = true
}

func (t *Trie) Search(word string) bool {
	current := t.Root

	for _, char := range word {
		if _, exists := current.Children[char]; !exists {
			return false
		}

		current = current.Children[char]
	}

	return current.IsEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	current := t.Root

	for _, char := range prefix {
		if _, exists := current.Children[char]; !exists {
			return false
		}

		current = current.Children[char]
	}
	return true
}

func (t *Trie) Delete(word string) bool {
	deleted, _ := deleteHelper(t.Root, word, 0)
	return deleted
}

func deleteHelper(node *Node, word string, depth int) (bool, bool) {
	if node == nil {
		return false, false
	}

	if depth == len(word) {
		if !node.IsEnd {
			return false, false
		}
		node.IsEnd = false
		// Delete only if it has no children
		return true, len(node.Children) == 0
	}

	if node.Children == nil {
		return false, false
	}

	ch := rune(word[depth])
	child, exists := node.Children[ch]
	if !exists || child == nil {
		return false, false
	}

	deleted, shouldDelete := deleteHelper(child, word, depth+1)

	if shouldDelete {
		delete(node.Children, ch)
	}

	// Return deletion status, and whether this node can be pruned
	return deleted, !node.IsEnd && len(node.Children) == 0
}

func (t *Trie) ListWords() []string {
	var words []string
	var dfs func(node *Node, prefix string)

	dfs = func(node *Node, prefix string) {
		if node.IsEnd {
			words = append(words, prefix)
		}

		for ch, child := range node.Children {
			dfs(child, prefix+string(ch))
		}
	}
	dfs(t.Root, "")

	return words
}
