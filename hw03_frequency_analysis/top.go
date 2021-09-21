package hw03frequencyanalysis

import (
	"bufio"
	"sort"
	"strings"
)

const topAmount = 10

type wordEntry struct {
	count int
	word  string
}

var (
	charsToTrim = "!?,.():"
	notWords    = map[string]bool{
		"":  true,
		"-": true,
		"*": true,
	}
)

func canonicalWord(word string) (canonical string) {
	return strings.Trim(strings.ToLower(word), charsToTrim)
}

func getWordsFrequency(input string) (frequency map[string]int) {
	frequency = make(map[string]int)

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := canonicalWord(scanner.Text())
		if !notWords[word] {
			frequency[word]++
		}
	}
	return frequency
}

func convertHistogram2Top(histogram map[string]int) (top []string) {
	entries := make([]wordEntry, 0, topAmount)
	for word, count := range histogram {
		entries = append(entries, wordEntry{count: count, word: word})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count == entries[j].count {
			return entries[i].word < entries[j].word
		}
		return entries[i].count > entries[j].count
	})

	top = make([]string, 0, topAmount)
	for i, entry := range entries {
		if i >= topAmount {
			break
		}
		top = append(top, entry.word)
	}

	return top
}

func Top10(input string) (top []string) {
	return convertHistogram2Top(getWordsFrequency(input))
}
