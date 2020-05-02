package main

import (
	"fmt"

	"github.com/wallberg/sandbox/taocp"
)

func main() {
	trie := taocp.NewPrefixTrie(3)
	trie.Add("abc")
	trie.Add("abe")
	trie.Add("ace")
	trie.Add("fun")
	trie.Add("gol")
	trie.Add("aaa")
	trie.Add("got")

	fmt.Println(trie.Size, trie.Count)

	words := make(chan string)
	go trie.Traverse(words)
	for word := range words {
		fmt.Println(word)
	}

}
