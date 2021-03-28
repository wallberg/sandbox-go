package sgb

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/packr"
)

// LoadWords loads the Stanford GraphBase 5-letter words, sorted by
// "commonality"
func LoadWords() ([]string, error) {
	// Load in ../taocp/assets/sgb-words.txt
	box := packr.NewBox("../taocp/assets")

	wordsString, err := box.FindString("sgb-words.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading assets/sgb-words.txt: %s", err)
	}

	words := strings.Split(wordsString, "\n")

	return words[0 : len(words)-1], nil
}

// LoadOSPD4 loads the words of The Official Scrabble Players Dictionary,
// version 4. If size is 0 then load all words, otherwise only load words of
// length size.
func LoadOSPD4(size int) ([]string, error) {
	// Load in ../taocp/assets/ospd4.txt
	box := packr.NewBox("../taocp/assets")

	wordsString, err := box.FindString("wordlists-ospd4.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading assets/wordlists-ospd4.txt: %s", err)
	}

	words := strings.Split(wordsString, "\n")

	words = words[0 : len(words)-1]

	if size == 0 {
		return words, nil
	}

	sizeWords := make([]string, 0)
	for _, word := range words {
		if len(word) == size {
			sizeWords = append(sizeWords, word)
		}
	}

	return sizeWords, nil
}
