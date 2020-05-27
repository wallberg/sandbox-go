package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wallberg/sandbox/taocp"
)

func printUsageAndExit(message string, code int) {
	if message != "" {
		fmt.Println(message)
		fmt.Println("")
	}

	fmt.Print(`USAGE
  toacp - Algorithms from The Art of Computer Programming

  taocp [--help] [-h] <command> ...

SUBCOMMANDS
    wr       Word Rectangles (m x n)
             -m, --m:       columns of words of size m; default is 5
             -n, --n:       rows of words of size n; default is 6
             -t, --threads: number of execution threads; default is 1
             -l, --limit:   limit to number of results per thread; default is 0, unlimited
`)

	os.Exit(code)
}

// cmdWr implements the 'wr' command
func cmdWr(m int, n int, threads int, limit int) {
	var trie taocp.Trie

	mTrie := taocp.NewCPrefixTrie(m)
	trie = &mTrie
	if m == 5 {
		taocp.LoadSGBWords(&trie)
	} else {
		taocp.LoadOSPD4Words(&trie, m)
	}

	nTrie := taocp.NewPrefixTrie(n)
	trie = &nTrie
	if n == 5 {
		taocp.LoadSGBWords(&trie)
	} else {
		taocp.LoadOSPD4Words(&trie, n)
	}

	words := make(chan string)
	go taocp.MultiWordRectangles(&mTrie, &nTrie, words, limit, threads)
	for word := range words {
		fmt.Println(word)
	}
}

func main() {

	wrCmd := flag.NewFlagSet("wr", flag.ExitOnError)
	wrM := wrCmd.Int("m", 5, "m")
	wrN := wrCmd.Int("n", 6, "n")
	wrThreads := wrCmd.Int("t", 1, "threads")
	wrLimit := wrCmd.Int("l", 0, "limit")

	if len(os.Args) < 2 {
		printUsageAndExit("expected subcommand: wr", 1)
	}

	switch os.Args[1] {

	case "wr":
		wrCmd.Parse(os.Args[2:])
		cmdWr(*wrM, *wrN, *wrThreads, *wrLimit)

	default:
		printUsageAndExit("expected subcommand: wr", 1)
	}
}
