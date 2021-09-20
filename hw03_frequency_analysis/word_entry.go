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
