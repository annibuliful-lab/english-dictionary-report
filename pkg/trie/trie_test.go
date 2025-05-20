package trie

import (
	"testing"
)

func TestTrieOperations(t *testing.T) {
	tr := NewTrie()

	tr.Insert("hello")
	tr.Insert("helium")
	tr.Insert("hero")

	// Test Search
	if !tr.Search("hello") {
		t.Error("Expected to find 'hello'")
	}
	if tr.Search("hell") {
		t.Error("Did not expect to find 'hell'")
	}

	// Test StartsWith
	if !tr.StartsWith("he") {
		t.Error("Expected prefix 'he'")
	}
	if tr.StartsWith("hex") {
		t.Error("Did not expect prefix 'hex'")
	}

	// Test Delete
	if !tr.Delete("hero") {
		t.Error("Expected 'hero' to be deleted")
	}
	if tr.Search("hero") {
		t.Error("Expected 'hero' to be removed")
	}
	if tr.Delete("herox") {
		t.Error("Expected 'herox' to fail deletion")
	}

	// Test ListWords
	words := tr.ListWords()
	expected := map[string]bool{"hello": true, "helium": true}
	for _, word := range words {
		if !expected[word] {
			t.Errorf("Unexpected word found: %s", word)
		}
		delete(expected, word)
	}
	for word := range expected {
		t.Errorf("Expected word missing: %s", word)
	}
}
