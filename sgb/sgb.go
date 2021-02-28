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
		return nil, fmt.Errorf("Error reading assets/sgb-words.txt: %s", err)
	}

	words := strings.Split(wordsString, "\n")

	return words[0 : len(words)-1], nil
}
