package report

import (
	"log"
	stringtools "myproject/stringtool"
	"sort"
	"strconv"
	"strings"
)

func BuildReport(lines []string) string {

	var b strings.Builder
	resultMap := make(map[string]int, len(lines))
	for _, line := range lines {
		before, after, found := strings.Cut(line, ":")
		if !found {
			log.Printf("пропущено: %q — нет разделителя", line)
			continue
		}
		after = strings.TrimSpace(after)
		value, err := strconv.Atoi(after)
		if err != nil {
			log.Printf("пропущено: %q — не число", line)
			continue
		}
		key := stringtools.ToLower(strings.TrimSpace(before))
		if key == "" {
			log.Printf("пропущено: %q — пустой ключ", line)
			continue
		}
		resultMap[key] = resultMap[key] + value
	}

	keys := make([]string, 0, len(resultMap))
	for k := range resultMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	header := "Отчёт по данным (" + strconv.Itoa(len(keys)) + " записей)\n"
	b.WriteString(header)
	var total int

	// Тело: цикл по ключам, много итераций — Builder, чтобы не было O(n²) от '+='.
	for _, k := range keys {
		value := strconv.Itoa(resultMap[k])
		b.WriteString("  ")
		b.WriteString(k)
		b.WriteString(" | ")
		b.WriteString(value)
		b.WriteString("\n")
		total = total + resultMap[k]
	}

	b.WriteString("Итого: " + strconv.Itoa(total))
	return b.String()
}
