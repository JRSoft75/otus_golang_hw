package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Регулярное выражение для удаления знаков препинания по краям слова.
var re = regexp.MustCompile(`^\p{P}+|\p{P}+$`)

func Top10(inputString string) []string {
	// Разделяем текст на слова
	words := strings.Fields(inputString)

	// Создаем карту для подсчета частоты слов
	wordCount := make(map[string]int)
	for _, word := range words {
		// Приводим слово к нижнему регистру
		word = strings.ToLower(word)
		// Удаляем знаки препинания по краям, если слово содержит буквы или цифры
		if containsLetterOrDigit(word) {
			word = re.ReplaceAllString(word, "")
		}
		// Пропускаем пустые строки и одиночные тире
		if word == "" || word == "-" {
			continue
		}
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

// Функция для проверки, содержит ли строка буквы или цифры.
func containsLetterOrDigit(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
	}
	return false
}
