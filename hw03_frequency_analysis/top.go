package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	wordsMap := make(map[string]int)
	for _, word := range strings.Fields(str) {
		wordsMap[word]++
	}

	sorted := rankByWordCount(wordsMap)
	maxCount := 10
	ans := []string{}

	for i, pair := range sorted {
		if i >= maxCount {
			break
		}
		ans = append(ans, pair.Key)
	}

	return ans
}

func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(pl)
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool {
	if p[i].Value > p[j].Value {
		return true
	} else if p[i].Value == p[j].Value {
		return p[i].Key < p[j].Key
	}
	return false
}
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
