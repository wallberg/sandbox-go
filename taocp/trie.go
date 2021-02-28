package taocp

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/wallberg/sandbox/sgb"
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

// CPrefixTrie represents a trie for words with the compressed prefix
// path by letter (a-z only). Compression comes in the form of storing the
// list of letters for a node in a linked list.  This compression reduces
// amount of memory necessary to store a node and makes stored letter
// traversal faster, at the expense of increased time to add a word.
type CPrefixTrie struct {
	Size  int    // fixed size of words in the trie
	Count int    // number of words in the trie
	Nodes []Link // the trie
}

// Link is a link in a singly linked list of letters.  The final link in
// the chain contains a right value of nil
type Link struct {
	Letter byte
	Node   int
	Right  *Link
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

// NewCPrefixTrie creates a new empty CPrefixTrie for words of length size
func NewCPrefixTrie(size int) CPrefixTrie {
	return CPrefixTrie{Size: size, Nodes: make([]Link, 0, 10)}
}

// Add adds a new word to the trie
func (trie *CPrefixTrie) Add(word string) {

	word = strings.ToLower(word)

	node := 0 // index into trie.Nodes)

	for i := 0; i < trie.Size; i++ {

		// Store letters a-z mapped to values 0..25
		letter := byte(word[i]) - 97

		// Add a new node, if necessary
		if node == len(trie.Nodes) {
			trie.Nodes = append(trie.Nodes, Link{0, 0, nil})
		}

		// Search the linked list to either find the existing entry for this
		// letter or insert a new one
		link := &trie.Nodes[node]
		for {
			if link.Right == nil || link.Letter > letter {
				// Insert here
				newLink := Link{link.Letter, link.Node, link.Right}
				link.Right = &newLink
				link.Letter = letter

				// Create new node, if necessary
				if i < trie.Size-1 {
					node = len(trie.Nodes)
				} else {
					node = FinalLetter
					trie.Count++
				}
				link.Node = node
				break

			} else if link.Letter == letter {
				// Letter already exists
				node = link.Node
				break
			}

			// Advance to next link
			link = link.Right
		}
	}
}

// Traverse sends to the out channel all words of the trie, in lexicographic
// order
func (trie *CPrefixTrie) Traverse(out chan string) {
	// Close the output channel on exit
	defer close(out)

	// node pointers, one per letter in the word
	link := make([]*Link, trie.Size)
	word := make([]byte, trie.Size)

	i := 0 // index of letter in the word, node, and letter arrays
	link[i] = &trie.Nodes[0]

	for {
		switch {

		case i < 0:
			// Traversal complete
			return

		case link[i].Right == nil:
			// Finished looking at all letters for this node
			i--
			if i >= 0 {
				link[i] = link[i].Right
			}

		default:
			// Assign letter to the word
			word[i] = link[i].Letter + 97
			if i == trie.Size-1 {
				// Visit the complete word
				out <- string(word)
			}

			if i < trie.Size-1 {
				// Advance to next node
				i++
				link[i] = &trie.Nodes[link[i-1].Node]
			} else {
				// Advance to next letter
				link[i] = link[i].Right
			}
		}
	}
}

// LoadSGBWords loads the Stanford GraphBase 5-letter words into a Trie
func LoadSGBWords(trie *Trie) error {
	words, err := sgb.LoadWords()
	if err != nil {
		return fmt.Errorf("Error reading assets/sgb-words.txt: %s", err)
	}

	for _, word := range words {
		(*trie).Add(word)
	}

	return nil
}

// LoadOSPD4Words loads the Official Scrabble Player's Dictionary, Version 4,
// n-letter words into a Trie
func LoadOSPD4Words(trie *Trie, n int) error {
	// Load in ./assets/.txt
	box := packr.NewBox("./assets")

	wordsString, err := box.FindString("ospd4.txt")
	if err != nil {
		return fmt.Errorf("Error reading assets/ospd4.txt: %s", err)
	}

	// Add each n-letter word to the Trie
	words := strings.Split(wordsString, "\n")
	for _, word := range words[0 : len(words)-1] {
		if len(word) == n {
			(*trie).Add(word)
		}
	}

	return nil
}
