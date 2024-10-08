package taocp

import (
	"testing"
)

func TestWordRectangles2by3(t *testing.T) {

	m := 2
	n := 3

	// 2 x 3 simple

	mTrie := NewCPrefixTrie(m)
	mTrie.Add("ab")
	//mTrie.Add("cd")
	mTrie.Add("ef")
	mTrie.Add("ag")
	mTrie.Add("ah")
	mTrie.Add("ai")
	mTrie.Add("ej")
	mTrie.Add("ek")

	nTrie := NewPrefixTrie(n)
	nTrie.Add("ace")
	nTrie.Add("bdf")
	nTrie.Add("alm")
	nTrie.Add("acn")
	nTrie.Add("bop")
	nTrie.Add("qrs")

	count := 0
	for range WordRectangles(&mTrie, &nTrie, 0, nil) {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 results; got %d", count)
	}

	mTrie.Add("cd")

	count = 0
	for result := range WordRectangles(&mTrie, &nTrie, 0, nil) {
		count++

		if count == 1 {
			if result != "ab:cd:ef" {
				t.Errorf("Expect result 1 to be ab:cd:ef, got %s", result)
			}
		}
	}

	if count != 1 {
		t.Errorf("Expected 1 result; got %d", count)
	}

}

func TestWordRectangles5x6(t *testing.T) {

	var trie Trie

	mTrie := NewCPrefixTrie(5)
	trie = &mTrie
	LoadSGBWords(&trie)

	nTrie := NewPrefixTrie(6)
	trie = &nTrie
	LoadOSPD4Words(&trie, 6)

	t.Run("SingleThread", func(t *testing.T) {
		singleWordRectangles5x6(t, &mTrie, &nTrie)
	})

	t.Run("MultiThread", func(t *testing.T) {
		multiWordRectangles(t, &mTrie, &nTrie)
	})
}

func singleWordRectangles5x6(t *testing.T, mTrie *CPrefixTrie, nTrie *PrefixTrie) {

	count := 0
	for result := range WordRectangles(mTrie, nTrie, 200, nil) {
		count++

		expected := ""
		if count == 1 {
			expected = "aargh:blare:lapin:atilt:tense:edged"
		} else if count == 191 {
			expected = "abaca:baths:bites:elude:sines:seers"
		}

		if expected != "" && result != expected {
			t.Errorf("Expected result for count %d to be %s; got %s", count, expected, result)
		}
	}

	if count != 200 {
		t.Errorf("Expected 200 results; got %d", count)
	}
}

func multiWordRectangles(t *testing.T, mTrie *CPrefixTrie, nTrie *PrefixTrie) {

	count := 0
	for range MultiWordRectangles(mTrie, nTrie, 5, 26, 0) {
		count++
	}

	if count != 130 {
		t.Errorf("Expected 130 results; got %d", count)
	}
}
