package maptools

import "sort"

type Student struct {
	Name    string
	Year    int
	Faculty string
}

func Frequency(words []string) map[string]int {
	result := make(map[string]int)

	for _, word := range words {
		result[word]++
	}
	return result
}

type WordCount struct {
	Word  string
	Count int
}

func FrequencyNoMap(words []string) []WordCount {
	result := []WordCount{}

	for i, word := range words {
		seen := false
		for j := 0; j < i; j++ {
			if words[j] == word {
				seen = true
				break
			}
		}
		if seen {
			continue
		}

		count := 0
		for _, w := range words {
			if w == word {
				count++
			}
		}
		result = append(result, WordCount{Word: word, Count: count})
	}

	return result
}

func GroupBy(src []Student) map[string]map[int][]Student {
	result := make(map[string]map[int][]Student)
	for _, student := range src {
		inner, ok := result[student.Faculty]
		if !ok {
			inner = make(map[int][]Student)
			result[student.Faculty] = inner
		}
		inner[student.Year] = append(inner[student.Year], student)
	}
	return result
}

func TopN(freq map[string]int, n int) []WordCount {
	if n <= 0 || 0 == len(freq) || freq == nil {
		return []WordCount{}
	}
	result := make([]WordCount, 0, len(freq))
	for word, count := range freq {
		result = append(result, WordCount{Word: word, Count: count})
	}

	if n >= len(freq) {
		return result
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Word < result[j].Word
	})
	if n < len(result) {
		result = result[:n]
	}
	return result
}
