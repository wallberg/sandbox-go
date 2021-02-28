package sgb

import "testing"

func TestLoadWords(t *testing.T) {
	words, err := LoadWords()

	if err != nil {
		t.Errorf("Error loading words: %v", err)
		return
	}

	if len(words) != 5757 {
		t.Errorf("Want 5757 words; got %d", len(words))
	}

	if words[0] != "which" {
		t.Errorf("Want first word of 'which'; got %s", words[0])
	}

	if words[len(words)-1] != "pupal" {
		t.Errorf("Want last word of 'pupal'; got %s", words[len(words)-1])
	}

}
