package taocp

import (
	"strings"
)

// Trie represents a trie for words of all the same size
type Trie interface {
	Add(string)
	Traverse(chan string)
}

// PrefixTrie represents a trie for words with the full prefix path
// by letter (a-z only)
type PrefixTrie struct {
	Size  int     // fixed size of words in the trie
	Count int     // number of words in the trie
	Nodes [][]int // the trie
}

// FinalLetter is the node value stored in the last letter of the word, instead
// of the next node in the sequence
const FinalLetter int = -1

// NewPrefixTrie creates a new empty PrefixTrie for words of length size
func NewPrefixTrie(size int) PrefixTrie {
	return PrefixTrie{Size: size, Nodes: make([][]int, 0, 10)}
}

// Add adds a new word to the trie
func (trie *PrefixTrie) Add(word string) {

	word = strings.ToLower(word)

	node := 0 // index into trie.Nodes)

	for i := 0; i < trie.Size; i++ {

		// Store letters a-z mapped to values 0..25
		l := byte(word[i]) - 97

		// Add a new node, if necessary
		if node == len(trie.Nodes) {
			trie.Nodes = append(trie.Nodes, make([]int, 26))
		}

		// Get value of next node
		switch nextNode := trie.Nodes[node][l]; {

		case nextNode == FinalLetter:
			// we are at the last letter and this word is already in the trie
			return

		case nextNode == 0:
			// this letter is not currently set
			if i < trie.Size-1 {
				// point to new node, to be created
				nextNode = len(trie.Nodes)
				trie.Nodes[node][l] = nextNode
				node = nextNode
			} else {
				// final node
				trie.Nodes[node][l] = FinalLetter
			}

		default:
			// this letter is already set, follow the node linek
			node = nextNode
		}
	}

	// Increment the word count
	trie.Count++
}

// Traverse sends to the out channel all words of the trie, in lexicographic
// order
func (trie *PrefixTrie) Traverse(out chan string) {
	// Close the output channel on exit
	defer close(out)

	// node pointers, one per letter in the word
	node := make([]int, trie.Size)
	letter := make([]byte, trie.Size)
	word := make([]byte, trie.Size)

	i := 0 // index of letter in the word, node, and letter arrays
	letter[i] = 0
	node[i] = 0

	for {
		switch {

		case i < 0:
			// Traversal complete
			return

		case letter[i] == 26:
			// Finished looking at all letters for this node
			i--
			if i >= 0 {
				letter[i]++
			}

		case trie.Nodes[node[i]][letter[i]] == 0:
			// Advance to next letter
			letter[i]++

		default:
			// Assign letter to the word
			word[i] = letter[i] + 97
			if i == trie.Size-1 {
				// Visit the complete word
				out <- string(word)
			}

			if i < trie.Size-1 {
				// Advance to next node
				i++
				node[i] = trie.Nodes[node[i-1]][letter[i-1]]
				letter[i] = 0
			} else {
				// Advance to next letter
				letter[i]++
			}
		}
	}
}
