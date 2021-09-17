package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const topAmount = 10

var (
	wordExtractor        = regexp.MustCompile(`[\S]+[-]?[\S]+|[^\s-]`)
	wordInsensitiveChars = regexp.MustCompile(`[!?,.*%@$^()]`)
)

type wordEntry struct {
	count int
	word  string
}

type wordEntrySorter struct {
	wordEntries []wordEntry
	by          func(w1, w2 *wordEntry) bool
}

func (s *wordEntrySorter) Len() int {
	return len(s.wordEntries)
}

func (s *wordEntrySorter) Swap(i, j int) {
	s.wordEntries[i], s.wordEntries[j] = s.wordEntries[j], s.wordEntries[i]
}

func (s *wordEntrySorter) Less(i, j int) bool {
	return s.by(&s.wordEntries[i], &s.wordEntries[j])
}

func canonicalWord(word string) (canonical string) {
	return wordInsensitiveChars.ReplaceAllString(strings.ToLower(word), "")
}

func entriesComparer(w1, w2 *wordEntry) bool {
	if w1.count == w2.count {
		return w1.word < w2.word
	}
	return w1.count > w2.count
}

func writeMatches2Channel(wordsChannel chan<- string, matches [][]string) {
	defer close(wordsChannel)

	for _, word := range matches {
		wordsChannel <- word[0]
	}
}

func readWords(wordsChannel <-chan string) map[string]int {
	frequency := make(map[string]int)
	for {
		word, opened := <-wordsChannel
		if !opened {
			return frequency
		}
		frequency[canonicalWord(word)]++
	}
}

func convertHistogram2Top(histogram map[string]int) (top []string) {
	entries := make([]wordEntry, 0, topAmount)
	for word, count := range histogram {
		entries = append(entries, wordEntry{count: count, word: word})
	}

	sort.Sort(&wordEntrySorter{wordEntries: entries, by: entriesComparer})

	limit := topAmount
	if len(entries) < topAmount {
		limit = len(entries)
	}

	top = make([]string, 0, topAmount)
	for _, entry := range entries[0:limit] {
		top = append(top, entry.word)
	}

	return top
}

func Top10(input string) (top []string) {
	wordsChannel := make(chan string)

	go writeMatches2Channel(wordsChannel, wordExtractor.FindAllStringSubmatch(input, -1))
	return convertHistogram2Top(readWords(wordsChannel))
}
