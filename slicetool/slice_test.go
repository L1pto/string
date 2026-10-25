package slicetools

import (
	"reflect"
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		name string
		src  []int
		pred func(int) bool
		want []int
	}{
		{"дубликаты", []int{1, 2, 3, 4, 5}, func(n int) bool { return n > 2 }, []int{3, 4, 5}},
		{"пустой", []int{}, func(n int) bool { return true }, []int{}},
		{"predicate false", []int{1, 2, 3, 4, 5}, func(n int) bool { return n > 6 }, []int{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Filter(c.src, c.pred)
			if !slices.Equal(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestFilterDoesNotMutateSource(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}

	sourceCopy := make([]int, len(src))

	copy(sourceCopy, src)

	_ = Filter(src, func(n int) bool { return true })

	if !slices.Equal(src, sourceCopy) {
		t.Fatalf("Исходный слайс изменился! Было: %v, Стало: %v", sourceCopy, src)
	}
}

func TestSquare(t *testing.T) {
	cases := []struct {
		name string
		src  []int
		want []int
	}{
		{"обычный случай", []int{1, 2, 3}, []int{1, 4, 9}},
		{"пустой", []int{}, []int{}},
		{"negative numbers", []int{-2, 3}, []int{4, 9}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Square(c.src)
			if !slices.Equal(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestDedup(t *testing.T) {
	cases := []struct {
		name string
		src  []int
		want []int
	}{
		{"обычный случай", []int{1, 3, 2, 3, 5, 1}, []int{1, 3, 2, 5}},
		{"пустой", []int{}, []int{}},
		{"negative numbers", []int{1, 2, -1, -2, -1, 0, 0}, []int{1, 2, -1, -2, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Dedup(c.src)
			if !slices.Equal(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestClone(t *testing.T) {
	cases := []struct {
		name string
		src  []int
		want []int
	}{
		{"обычный случай", []int{1, 2, 3}, []int{1, 2, 3}},
		{"пустой", []int{}, []int{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Clone(c.src)
			if !slices.Equal(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestCloneDoesNotMutateSource(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}

	sourceCopy := make([]int, len(src))

	copy(sourceCopy, src)

	got := Clone(src)
	got[0] = 777
	if !slices.Equal(src, sourceCopy) {
		t.Fatalf("Исходный слайс изменился! Было: %v, Стало: %v", sourceCopy, src)
	}
}

func TestChunk(t *testing.T) {
	cases := []struct {
		name string
		src  []int
		size int
		want [][]int
	}{
		{"обычный случай", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{"пустой", []int{}, 1, nil},
		{"size > len(src)", []int{1, 2, 3}, 10, [][]int{{1, 2, 3}}},
		{"size == len(src)", []int{1, 2, 3}, 3, [][]int{{1, 2, 3}}},
		{"size == 1", []int{1, 2, 3}, 3, [][]int{{1, 2, 3}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Chunk(c.src, c.size)
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestChunkPanicsOnZeroSize(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("Chunk не паникует а должен")
		}

		if msg, ok := r.(string); !ok || msg != "slicetools: size must be > 0" {
			t.Fatalf("Неверное сообщение об ошибке")
		}
	}()
	Chunk([]int{1, 2, 3}, 0)
}

func TestChunkDoesNotMutateSource(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}

	sourceCopy := make([]int, len(src))

	copy(sourceCopy, src)

	chunks := Chunk(src, 3)
	chunks[0][0] = 999
	if !slices.Equal(src, sourceCopy) {
		t.Fatalf("Исходный слайс изменился! Было: %v, Стало: %v", sourceCopy, src)
	}
}
