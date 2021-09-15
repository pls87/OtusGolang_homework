package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const topAmount int = 10

var (
	separator    = regexp.MustCompile(`\s+`)
	wordsChannel = make(chan string)
	frequency    = make(map[string]int)
)

type wordEntry struct {
	word  string
	count int
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

func countAndLexicalComparer(w1, w2 *wordEntry) bool {
	if w1.count == w2.count {
		return w1.word < w2.word
	}
	return w1.count > w2.count
}

func writeWords2Channel(words []string) {
	defer close(wordsChannel)

	for _, word := range words {
		wordsChannel <- word
	}
}

func readChannel() {
	for {
		word, opened := <-wordsChannel
		if !opened {
			break
		}
		frequency[word]++
	}
}

func Top10(input string) (top10 []string) {
	input = strings.Trim(input, " ")
	if input == "" {
		return make([]string, 0, 0)
	}

	words := separator.Split(input, -1)

	go writeWords2Channel(words)
	readChannel()

	entries := make([]wordEntry, 0, topAmount)
	for word, count := range frequency {
		entries = append(entries, wordEntry{word, count})
	}

	sort.Sort(&wordEntrySorter{entries, countAndLexicalComparer})

	top10 = make([]string, 0, topAmount)
	for _, value := range entries[0:topAmount] {
		top10 = append(top10, value.word)
	}
	return top10
}
