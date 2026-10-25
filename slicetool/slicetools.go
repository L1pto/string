package slicetools

func Filter(src []int, pred func(int) bool) []int {
	result := make([]int, 0, len(src))
	for _, v := range src {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

func Square(src []int) []int {
	result := make([]int, len(src))
	for i := range src {
		result[i] = src[i] * src[i]
	}
	return result
}

func Dedup(src []int) []int {
	seen := map[int]bool{}
	result := make([]int, 0, len(src))
	for _, v := range src {
		if seen[v] {
			continue
		}
		seen[v] = true
		result = append(result, v)
	}
	return result
}

func Clone(src []int) []int {
	result := make([]int, len(src))
	copy(result, src)
	return result
}

func Chunk(src []int, size int) [][]int {

	if size <= 0 {
		panic("slicetools: size must be > 0")
	}

	if len(src) == 0 {
		return nil
	}

	n := (len(src) + size - 1) / size
	chunks := make([][]int, 0, n)
	for f := 0; f < n; f++ {
		i := f * size
		j := min((f+1)*size, len(src))
		innerChunk := make([]int, j-i)
		copy(innerChunk, src[i:j])
		chunks = append(chunks, innerChunk)
	}
	return chunks
}
