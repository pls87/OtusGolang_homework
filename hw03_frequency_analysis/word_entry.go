package hw03frequencyanalysis

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

func entriesComparer(w1, w2 *wordEntry) bool {
	if w1.count == w2.count {
		return w1.word < w2.word
	}
	return w1.count > w2.count
}
