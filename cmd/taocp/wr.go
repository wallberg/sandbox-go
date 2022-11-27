package main

import (
	"fmt"
	"log"

	"github.com/wallberg/sandbox-go/taocp"
)

// initialize this command by adding it to the parser
func init() {
	var command wrCommand

	_, err := parser.AddCommand("wr",
		"Word Rectangles",
		"Generate m x n Word Rectangles",
		&command,
	)
	if err != nil {
		log.Fatalf("Error adding wr command: %v", err)
	}
}

type wrCommand struct {
	M       int `short:"m" long:"m" description:"number of rows" default:"5"`
	N       int `short:"n" long:"n" description:"number of columngs" default:"6"`
	Threads int `short:"t" long:"threads" description:"number of execution threads" default:"1"`
	Index   int `short:"i" long:"index" description:"thread to run (0 means all)" default:"0"`
	Limit   int `short:"l" long:"limit" description:"number of results per thread (0 means all)" default:"0"`
}

func (command wrCommand) Execute(args []string) error {

	var trie taocp.Trie

	mTrie := taocp.NewCPrefixTrie(command.M)
	trie = &mTrie
	if command.M == 5 {
		taocp.LoadSGBWords(&trie)
	} else {
		taocp.LoadOSPD4Words(&trie, command.M)
	}

	nTrie := taocp.NewPrefixTrie(command.N)
	trie = &nTrie
	if command.N == 5 {
		taocp.LoadSGBWords(&trie)
	} else {
		taocp.LoadOSPD4Words(&trie, command.N)
	}

	words := make(chan string)

	go taocp.MultiWordRectangles(&mTrie, &nTrie, words,
		command.Limit, command.Threads, command.Index)

	for word := range words {
		fmt.Println(word)
	}

	return nil
}
