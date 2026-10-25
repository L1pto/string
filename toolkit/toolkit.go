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
	// 1. очистка и нормализация
	clean := stringtools.CleanNonLetters(stringtools.ToLower(text))
	fmt.Println("clean:", clean)

	// 2. разбиение на слова — тут нужен SplitWords, см. ниже
	words := stringtools.SplitWords(clean)
	fmt.Println("words:", words)

	// 3. частота
	freq := maptools.Frequency(words)
	fmt.Println("freq:", freq)

	// 4. топ-3 частых слов (через slicetools)
	top := maptools.TopN(freq, 3) // использует slicetools.Filter / sort
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
