package toolkit

import (
	"fmt"
	maptools "myproject/maptool"
	"myproject/report"
	slicetools "myproject/slicetool"
	stringtools "myproject/stringtool"
)

func RunNumbersDemo() {
	nums := []int{3, 1, 3, 2, 1, 4, 2, 3}
	fmt.Println("src:", nums)
	fmt.Println("Filter n>2:", slicetools.Filter(nums, func(n int) bool { return n > 2 }))
	fmt.Println("Square:", slicetools.Square(nums))
	fmt.Println("Dedup:", slicetools.Dedup(nums))
	fmt.Println("Chunk(3):", slicetools.Chunk(nums, 3))
}

func RunTextDemo(text string) {
	// 1. Разбиваем исходный текст на слова
	rawWords := stringtools.SplitWords(stringtools.ToLower(text))

	// 2. Нормализуем каждый токен: оставляем только буквы
	words := make([]string, 0, len(rawWords))
	for _, w := range rawWords {
		cleaned := stringtools.CleanNonLetters(w)
		if cleaned != "" {
			words = append(words, cleaned)
		}
	}

	freq := maptools.Frequency(words)
	top := maptools.TopN(freq, 3)
	fmt.Println("freq:", freq)
	fmt.Println("top-3:", top)
}

func RunReportDemo() {
	lines := []string{
		"Яблоки: 10", "Бананы: 5", "яблоки: 3", "  Груши  :   7  ",
		"БезРазделителя", "Хлеб: неЧисло", "Кефир:3.5",
	}
	fmt.Println("src:", lines)
	fmt.Println(report.BuildReport(lines))
}
