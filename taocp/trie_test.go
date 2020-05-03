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

	cPrefixTrie := NewCPrefixTrie(1)
	trie = &cPrefixTrie
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

func TestNewCPrefixTrie(t *testing.T) {
	trie := NewCPrefixTrie(3)

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

func TestPrefixTrieAdd(t *testing.T) {
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

func TestCPrefixTrieAdd(t *testing.T) {
	trie := NewCPrefixTrie(3)

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

func TestPrefixTrieTraverse(t *testing.T) {
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

func TestCPrefixTrieTraverse(t *testing.T) {
	trie := NewCPrefixTrie(3)

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

func TestPrefixTrieLoadSGBWords(t *testing.T) {
	var trie Trie
	prefixTrie := NewPrefixTrie(5)
	trie = &prefixTrie
	err := LoadSGBWords(&trie)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if prefixTrie.Count != 5757 {
		t.Errorf("Expected trie.Count of 5757; got %d", prefixTrie.Count)
	}
}

func TestCPrefixTrieLoadSGBWords(t *testing.T) {
	var trie Trie
	cPrefixTrie := NewCPrefixTrie(5)
	trie = &cPrefixTrie
	err := LoadSGBWords(&trie)

	if err != nil {
		t.Errorf("Error: %s", err)
	}

	if cPrefixTrie.Count != 5757 {
		t.Errorf("Expected trie.Count of 5757; got %d", cPrefixTrie.Count)
	}

	words := make(chan string)
	go trie.Traverse(words)

	i := 0
	for word := range words {
		expected := ""
		switch i {
		case 0:
			expected = "aargh"
		case 1:
			expected = "abaca"
		case 2:
			expected = "abaci"
		case 428:
			expected = "berry"
		case 1248:
			expected = "deque"
		case 2968:
			expected = "mails"
		case 4458:
			expected = "skews"
		case 5754:
			expected = "zooks"
		case 5755:
			expected = "zooms"
		case 5756:
			expected = "zowie"
		}

		if expected != "" && expected != word {
			t.Errorf("Expected %s at position %d; got %s", expected, i, word)
		}

		i++
	}
}

func TestPrefixTrieLoadOSPD4Words(t *testing.T) {
	var trie Trie
	var prefixTrie PrefixTrie

	prefixTrie = NewPrefixTrie(6)
	trie = &prefixTrie
	err := LoadOSPD4Words(&trie, 6)

	if err != nil {
		t.Errorf("Error: %s", err)

	} else if prefixTrie.Count != 15727 {
		t.Errorf("Expected trie.Count of 15727; got %d", prefixTrie.Count)
	}

	prefixTrie = NewPrefixTrie(2)
	trie = &prefixTrie
	err = LoadOSPD4Words(&trie, 2)

	if err != nil {
		t.Errorf("Error: %s", err)

	} else if prefixTrie.Count != 101 {
		t.Errorf("Expected trie.Count of 101; got %d", prefixTrie.Count)
	}
}

func TestCPrefixTrieLoadOSPD4Words(t *testing.T) {
	var trie Trie
	var cPrefixTrie CPrefixTrie

	cPrefixTrie = NewCPrefixTrie(6)
	trie = &cPrefixTrie
	err := LoadOSPD4Words(&trie, 6)

	if err != nil {
		t.Errorf("Error: %s", err)

	} else if cPrefixTrie.Count != 15727 {
		t.Errorf("Expected trie.Count of 15727; got %d", cPrefixTrie.Count)
	}

	cPrefixTrie = NewCPrefixTrie(2)
	trie = &cPrefixTrie
	err = LoadOSPD4Words(&trie, 2)

	if err != nil {
		t.Errorf("Error: %s", err)

	} else if cPrefixTrie.Count != 101 {
		t.Errorf("Expected trie.Count of 101; got %d", cPrefixTrie.Count)
	}
}
