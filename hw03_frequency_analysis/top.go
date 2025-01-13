package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	words := strings.Fields(input)
	frequency := make(map[string]int)

	// считаем сколько раз встречаются слова в мапу => 'one': 1; 'two': 2; 'three': 2
	for _, word := range words {
		frequency[word]++
	}

	// запоминаем ключи == уникальные слова
	uniqueWords := make([]string, 0)
	for word := range frequency {
		uniqueWords = append(uniqueWords, word)
	}

	// сортируем уникальные слова
	sort.Slice(uniqueWords, func(i, j int) bool {
		leftWord := uniqueWords[i]
		rightWord := uniqueWords[j]
		if frequency[leftWord] == frequency[rightWord] { // если частота двух слов одинаковая, сортируем лексикографически
			return strings.Compare(leftWord, rightWord) < 0
		}
		return frequency[leftWord] > frequency[rightWord]
	})

	if len(uniqueWords) > 10 {
		return uniqueWords[:10]
	}

	return uniqueWords
}
