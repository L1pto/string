package report

import "testing"

func TestBuildReport(t *testing.T) {
	cases := []struct {
		name string
		src  []string
		want string
	}{
		{"empty", []string{}, "Отчёт по данным (0 записей)\nИтого: 0"},
		{"пустой", nil, "Отчёт по данным (0 записей)\nИтого: 0"},
		{"один", []string{"Хлеб: 5"}, "Отчёт по данным (1 записей)\n  хлеб | 5\nИтого: 5"},
		{"нормализация ключа", []string{"Яблоко: 10", "яблоко: 7"}, "Отчёт по данным (1 записей)\n  яблоко | 17\nИтого: 17"},
		{"пробелы вокруг", []string{"  Груши  :   7  "}, "Отчёт по данным (1 записей)\n  груши | 7\nИтого: 7"},
		{"отрицательное", []string{"Сыр: -2"}, "Отчёт по данным (1 записей)\n  сыр | -2\nИтого: -2"},
		{"сортировка", []string{"Б: 1", "А: 2", "В: 3"}, "Отчёт по данным (3 записей)\n  а | 2\n  б | 1\n  в | 3\nИтого: 6"},
		{"двоеточие в значении", []string{"ключ: 1:2"}, "Отчёт по данным (0 записей)\nИтого: 0"},
		{"плюс в значении", []string{"x: +5"}, "Отчёт по данным (1 записей)\n  x | 5\nИтого: 5"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := BuildReport(c.src)
			if got != c.want {
				t.Fatalf("unexpected output %s, %s", got, c.want)
			}
		})
	}
}

func TestBuildReportSkipsInvalidLines(t *testing.T) {
	lines := []string{
		"Молоко: 4",
		"БезРазделителя", // нет ':'
		"Хлеб: неЧисло",  // не число
		"Кефир:3.5",      // float → Atoi не парсит
		"Сыр: 12",
		"",       // пустая строка
		"   ",    // только пробелы
		": 5",    // пустой ключ
		"   : 7", // пустой ключ после TrimSpace
	}
	got := BuildReport(lines)
	want := "Отчёт по данным (2 записей)\n  молоко | 4\n  сыр | 12\nИтого: 16"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBuildReportSortsKeys(t *testing.T) {
	// Вход заведомо не в алфавитном порядке.
	// Если sort.Strings убрать — порядок станет случайным,
	// и тест будет падать нестабильно. Чтобы поймать регрессию
	// детерминированно, прогоняем несколько раз и требуем
	// одинаковый результат.
	src := []string{"Яблоко: 1", "Апельсин: 2", "Банан: 3", "Вишня: 4"}

	first := BuildReport(src)
	for i := 0; i < 20; i++ {
		if got := BuildReport(src); got != first {
			t.Fatalf("вывод недетерминирован на итерации %d:\nfirst=%q\ngot  =%q", i, first, got)
		}
	}

	want := "Отчёт по данным (4 записей)\n" +
		"  апельсин | 2\n" +
		"  банан | 3\n" +
		"  вишня | 4\n" +
		"  яблоко | 1\n" +
		"Итого: 10"
	if first != want {
		t.Fatalf("got %q, want %q", first, want)
	}
}
