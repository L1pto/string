package maptools

import (
	"maps"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestFrequency(t *testing.T) {
	cases := []struct {
		name string
		src  []string
		want map[string]int
	}{
		{"mixed", []string{"Ch", "Rus", "UK", "Rus", "Jp", "Ch"}, map[string]int{"Ch": 2, "Jp": 1, "Rus": 2, "UK": 1}},
		{"zero", []string{}, map[string]int{}},
		{"empty string", []string{"", "a", ""}, map[string]int{"": 2, "a": 1}},
		{"регистрозависимость", []string{"Go", "go", "GO"}, map[string]int{"Go": 1, "go": 1, "GO": 1}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Frequency(c.src)
			if !maps.Equal(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestFrequency_MissingKey(t *testing.T) {
	got := Frequency([]string{"a"})
	if _, ok := got["b"]; ok {
		t.Fatal("отсутствующий ключ не должен появляться")
	}
	if len(got) != 1 {
		t.Fatalf("ожидали 1 ключ, got %d: %v", len(got), got)
	}
}

func eqOuter(got, want map[string]map[int][]Student) bool {
	return maps.EqualFunc(got, want, eqInner)
}

func eqInner(got, want map[int][]Student) bool {
	return maps.EqualFunc(got, want, func(a, b []Student) bool {
		return slices.Equal(a, b) // Student сравним (все поля comparable)
	})
}

func TestGroupBy(t *testing.T) {
	// want должен быть отсортирован по Name внутри каждой группы
	cases := []struct {
		name string
		src  []Student
		want map[string]map[int][]Student
	}{
		{
			name: "mixed: несколько факультетов, разные годы",
			src: []Student{
				{Name: "Иванов Иван", Year: 2005, Faculty: "Математический"},
				{Name: "Петрова Анна", Year: 2006, Faculty: "Математический"},
				{Name: "Сидоров Пётр", Year: 2004, Faculty: "Математический"},
				{Name: "Кузнецова Мария", Year: 2003, Faculty: "Математический"},
				{Name: "Смирнов Алексей", Year: 2003, Faculty: "Математический"},
				{Name: "Морозов Дмитрий", Year: 2005, Faculty: "Физический"},
				{Name: "Соколов Артём", Year: 2005, Faculty: "Физический"},
				{Name: "Новиков Сергей", Year: 2006, Faculty: "Информационных технологий"},
				{Name: "Фёдорова Дарья", Year: 2005, Faculty: "Информационных технологий"},
			},
			want: map[string]map[int][]Student{
				"Математический": {
					2005: {{Name: "Иванов Иван", Year: 2005, Faculty: "Математический"}},
					2006: {{Name: "Петрова Анна", Year: 2006, Faculty: "Математический"}},
					2004: {{Name: "Сидоров Пётр", Year: 2004, Faculty: "Математический"}},
					2003: {
						{Name: "Кузнецова Мария", Year: 2003, Faculty: "Математический"},
						{Name: "Смирнов Алексей", Year: 2003, Faculty: "Математический"},
					},
				},
				"Физический": {
					2005: {
						{Name: "Морозов Дмитрий", Year: 2005, Faculty: "Физический"},
						{Name: "Соколов Артём", Year: 2005, Faculty: "Физический"},
					},
				},
				"Информационных технологий": {
					2006: {{Name: "Новиков Сергей", Year: 2006, Faculty: "Информационных технологий"}},
					2005: {{Name: "Фёдорова Дарья", Year: 2005, Faculty: "Информационных технологий"}},
				},
			},
		},
		{
			name: "zero: пустой срез",
			src:  []Student{},
			want: map[string]map[int][]Student{},
		},
		{
			name: "nil: nil-срез",
			src:  nil,
			want: map[string]map[int][]Student{},
		},
		{
			name: "разные факультеты, без пересечения по годам",
			src: []Student{
				{Name: "A", Year: 2001, Faculty: "Физ"},
				{Name: "B", Year: 2002, Faculty: "Хим"},
				{Name: "C", Year: 2003, Faculty: "Био"},
			},
			want: map[string]map[int][]Student{
				"Физ": {2001: {{Name: "A", Year: 2001, Faculty: "Физ"}}},
				"Хим": {2002: {{Name: "B", Year: 2002, Faculty: "Хим"}}},
				"Био": {2003: {{Name: "C", Year: 2003, Faculty: "Био"}}},
			},
		},
		{
			name: "один факультет, разные годы",
			src: []Student{
				{Name: "A", Year: 2001, Faculty: "Физ"},
				{Name: "B", Year: 2002, Faculty: "Физ"},
				{Name: "C", Year: 2003, Faculty: "Физ"},
			},
			want: map[string]map[int][]Student{
				"Физ": {
					2001: {{Name: "A", Year: 2001, Faculty: "Физ"}},
					2002: {{Name: "B", Year: 2002, Faculty: "Физ"}},
					2003: {{Name: "C", Year: 2003, Faculty: "Физ"}},
				},
			},
		},
		{
			name: "один факультет, один год, несколько студентов",
			src: []Student{
				{Name: "A", Year: 2001, Faculty: "Физ"},
				{Name: "B", Year: 2001, Faculty: "Физ"},
				{Name: "C", Year: 2001, Faculty: "Физ"},
			},
			want: map[string]map[int][]Student{
				"Физ": {
					2001: {
						{Name: "A", Year: 2001, Faculty: "Физ"},
						{Name: "B", Year: 2001, Faculty: "Физ"},
						{Name: "C", Year: 2001, Faculty: "Физ"},
					},
				},
			},
		},
		{
			name: "один студент",
			src:  []Student{{Name: "Solo", Year: 2000, Faculty: "Физ"}},
			want: map[string]map[int][]Student{
				"Физ": {2000: {{Name: "Solo", Year: 2000, Faculty: "Физ"}}},
			},
		},
		{
			name: "одинаковые имена, но разные годы/факультеты",
			src: []Student{
				{Name: "Иванов", Year: 2001, Faculty: "Физ"},
				{Name: "Иванов", Year: 2001, Faculty: "Хим"},
				{Name: "Иванов", Year: 2002, Faculty: "Физ"},
			},
			want: map[string]map[int][]Student{
				"Физ": {
					2001: {{Name: "Иванов", Year: 2001, Faculty: "Физ"}},
					2002: {{Name: "Иванов", Year: 2002, Faculty: "Физ"}},
				},
				"Хим": {
					2001: {{Name: "Иванов", Year: 2001, Faculty: "Хим"}},
				},
			},
		},
	}

	sortStudents := func(m map[string]map[int][]Student) {
		for _, years := range m {
			for _, students := range years {
				slices.SortFunc(students, func(a, b Student) int {
					return strings.Compare(a.Name, b.Name)
				})
			}
		}
	}
	for _, c := range cases {

		t.Run(c.name, func(t *testing.T) {
			got := GroupBy(c.src)
			sortStudents(got)

			if !eqOuter(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestGroupBy_Invariants(t *testing.T) {
	src := []Student{
		{Name: "A", Year: 2001, Faculty: "Физ"},
		{Name: "B", Year: 2001, Faculty: "Физ"},
		{Name: "C", Year: 2002, Faculty: "Хим"},
	}
	got := GroupBy(src)

	// 1) каждый студент попал ровно в одну группу
	total := 0
	for faculty, years := range got {
		for year, students := range years {
			for _, s := range students {
				if s.Faculty != faculty || s.Year != year {
					t.Errorf("студент %v попал в группу %q/%d", s, faculty, year)
				}
				total++
			}
		}
	}
	if total != len(src) {
		t.Errorf("в группах %d студентов, во входе %d", total, len(src))
	}

	// 2) пустых групп нет
	for faculty, years := range got {
		if len(years) == 0 {
			t.Errorf("факультет %q с пустой мапой годов", faculty)
		}
		for year, students := range years {
			if len(students) == 0 {
				t.Errorf("пустая группа %q/%d", faculty, year)
			}
		}
	}
	wantFaculties := map[string]bool{"Физ": true, "Хим": true}
	for faculty := range got {
		if !wantFaculties[faculty] {
			t.Errorf("лишний факультет %q", faculty)
		}
	}
}

func TestFrequencyNoMap(t *testing.T) {
	words := []string{"a", "b", "a", "c", "b", "a"}
	got := FrequencyNoMap(words)

	m := make(map[string]int, len(got))
	for _, wc := range got {
		m[wc.Word] = wc.Count
	}
	if !maps.Equal(m, Frequency(words)) {
		t.Fatalf("counts: got %v, want %v", m, Frequency(words))
	}

	gotOrder := make([]string, len(got))
	for i, wc := range got {
		gotOrder[i] = wc.Word
	}
	if !slices.Equal(gotOrder, []string{"a", "b", "c"}) {
		t.Fatalf("порядок: got %v, want [a b c]", gotOrder)
	}
	if got := FrequencyNoMap(nil); len(got) != 0 {
		t.Fatalf("для nil ожидали пустой результат, got %v", got)
	}
}

// Вспомогательная функция для удобного сравнения без учёта nil vs empty
func equalCounts(a, b []WordCount) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func copyFreq(m map[string]int) map[string]int {
	if m == nil {
		return nil
	}
	c := make(map[string]int, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}

func TestTopN(t *testing.T) {
	t.Run("regular case n=2", func(t *testing.T) {
		freq := map[string]int{"a": 3, "b": 1, "c": 2, "d": 5}
		// Отсортировано по убыванию Count: d(5), a(3), c(2), b(1)
		// top-2: d, a
		got := TopN(freq, 2)
		want := []WordCount{{Word: "d", Count: 5}, {Word: "a", Count: 3}}
		if !equalCounts(got, want) {
			t.Fatalf("TopN(%v, 2) = %#v, want %#v", freq, got, want)
		}
	})

	t.Run("tie-break by word ascending", func(t *testing.T) {
		// Все с одинаковым Count → порядок должен быть по слову (a, b, c, d)
		freq := map[string]int{"c": 5, "a": 5, "d": 5, "b": 5}
		got := TopN(freq, 3)
		want := []WordCount{
			{Word: "a", Count: 5},
			{Word: "b", Count: 5},
			{Word: "c", Count: 5},
		}
		if !equalCounts(got, want) {
			t.Fatalf("TopN(%v, 3) = %#v, want %#v", freq, got, want)
		}
	})

	t.Run("tie-break mixed counts", func(t *testing.T) {
		// Count: x=1, y=1, z=2, w=2, v=3
		// Ожидаемый порядок: v(3), w(2), z(2), x(1), y(1)
		freq := map[string]int{"x": 1, "y": 1, "z": 2, "w": 2, "v": 3}
		got := TopN(freq, 4)
		want := []WordCount{
			{Word: "v", Count: 3},
			{Word: "w", Count: 2},
			{Word: "z", Count: 2},
			{Word: "x", Count: 1},
		}
		if !equalCounts(got, want) {
			t.Fatalf("TopN(%v, 4) = %#v, want %#v", freq, got, want)
		}
	})

	t.Run("n > len(freq)", func(t *testing.T) {
		freq := map[string]int{"a": 1, "b": 2}
		got := TopN(freq, 10)
		want := []WordCount{{Word: "a", Count: 1}, {Word: "b", Count: 2}}
		if !equalCounts(got, want) {
			t.Fatalf("TopN(%v, 10) = %#v, want %#v", freq, got, want)
		}
	})

	t.Run("n == len(freq)", func(t *testing.T) {
		freq := map[string]int{"a": 1, "b": 2, "c": 3}
		got := TopN(freq, 3)
		// Порядок не гарантирован (sort не вызывается), поэтому сравниваем как множества
		if len(got) != 3 {
			t.Fatalf("expected 3 items, got %d: %#v", len(got), got)
		}
		seen := map[string]int{}
		for _, wc := range got {
			seen[wc.Word] = wc.Count
		}
		if !reflect.DeepEqual(seen, freq) {
			t.Fatalf("TopN(%v, 3) returned mismatched set: %#v", freq, got)
		}
	})

	t.Run("n == 0", func(t *testing.T) {
		freq := map[string]int{"a": 1, "b": 2}
		got := TopN(freq, 0)
		if len(got) != 0 {
			t.Fatalf("expected empty, got %#v", got)
		}
		if got == nil {
			t.Fatalf("expected non-nil empty slice")
		}
	})

	t.Run("n < 0", func(t *testing.T) {
		freq := map[string]int{"a": 1, "b": 2}
		got := TopN(freq, -5)
		if len(got) != 0 {
			t.Fatalf("expected empty, got %#v", got)
		}
		if got == nil {
			t.Fatalf("expected non-nil empty slice")
		}
	})

	t.Run("empty map", func(t *testing.T) {
		freq := map[string]int{}
		got := TopN(freq, 3)
		if len(got) != 0 {
			t.Fatalf("expected empty, got %#v", got)
		}
		if got == nil {
			t.Fatalf("expected non-nil empty slice")
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var freq map[string]int // nil
		got := TopN(freq, 3)
		if len(got) != 0 {
			t.Fatalf("expected empty, got %#v", got)
		}
		if got == nil {
			t.Fatalf("expected non-nil empty slice")
		}
	})

	t.Run("does not mutate input map", func(t *testing.T) {
		freq := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
		before := copyFreq(freq)
		_ = TopN(freq, 2)
		if !reflect.DeepEqual(freq, before) {
			t.Fatalf("input map was mutated: before=%v after=%v", before, freq)
		}
	})

	t.Run("single element n=1", func(t *testing.T) {
		freq := map[string]int{"a": 42}
		got := TopN(freq, 1)
		want := []WordCount{{Word: "a", Count: 42}}
		if !equalCounts(got, want) {
			t.Fatalf("TopN(%v, 1) = %#v, want %#v", freq, got, want)
		}
	})
}
