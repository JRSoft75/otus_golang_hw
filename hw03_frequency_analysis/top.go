package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(inputString string) []string {
	// Разделяем текст на слова
	words := strings.Fields(inputString)

	// Создаем карту для подсчета частоты слов
	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}

	// Создаем слайс для хранения слов и их частоты
	type wordFrequency struct {
		word  string
		count int
	}

	frequencies := make([]wordFrequency, 0, len(wordCount))
	for word, count := range wordCount {
		frequencies = append(frequencies, wordFrequency{word, count})
	}

	// Сортируем по частоте, а затем лексикографически
	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].count == frequencies[j].count {
			return frequencies[i].word < frequencies[j].word
		}
		return frequencies[i].count > frequencies[j].count
	})

	// Извлекаем первые 10 слов
	result := make([]string, 0, 10)
	for i := 0; i < len(frequencies) && i < 10; i++ {
		result = append(result, frequencies[i].word)
	}

	return result
}
