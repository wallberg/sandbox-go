package taocp

import (
	"testing"
)

func TestTrie(t *testing.T) {
	// Basically we succeed if there are no compile time errors
	var trie Trie
	prefixTrie := NewPrefixTrie(1)
	trie = &prefixTrie
	trie.Add("a")
}

func TestNewPrefixTrie(t *testing.T) {
	trie := NewPrefixTrie(3)

	if trie.Size != 3 {
		t.Errorf("Expected trie.Size value of 3; got %d", trie.Size)
	}

	if trie.Count != 0 {
		t.Errorf("Expected trie.Count value of 0; got %d", trie.Count)
	}

	if nodesLen := len(trie.Nodes); nodesLen != 0 {
		t.Errorf("Expected len(trie.Nodes) value of 0; got %d", nodesLen)
	}

}

func TestAdd(t *testing.T) {
	trie := NewPrefixTrie(3)

	trie.Add("abc")

	if trie.Count != 1 {
		t.Errorf("Expected trie.Count value of 1; got %d", trie.Count)
	}

	trie.Add("abc") // duplicate
	trie.Add("abe")
	trie.Add("ace")
	trie.Add("got")
	trie.Add("fun")
	trie.Add("gol")
	trie.Add("aaa")
	trie.Add("got") // duplicate

	if trie.Count != 7 {
		t.Errorf("Expected trie.Count value of 7; got %d", trie.Count)
	}

}

func TestTraverse(t *testing.T) {
	trie := NewPrefixTrie(3)

	trie.Add("abc")
	trie.Add("abe")
	trie.Add("ace")
	trie.Add("got")
	trie.Add("fun")
	trie.Add("gol")
	trie.Add("aaa")

	wordsChannel := make(chan string)
	words := make([]string, 0)
	go trie.Traverse(wordsChannel)
	for word := range wordsChannel {
		words = append(words, word)
	}

	expectedWords := []string{
		"aaa",
		"abc",
		"abe",
		"ace",
		"fun",
		"gol",
		"got",
	}

	match := false

	if len(words) == len(expectedWords) {

		match = true
		for i := 0; i < len(words); i++ {
			if words[i] != expectedWords[i] {
				match = false
				break
			}
		}
	}

	if !match {
		t.Errorf("Expected word array of %s; got %s", expectedWords, words)
	}
}
