package stringtools

import (
	"slices"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestStringModel(t *testing.T) {
	s := "Привет, 🌍!"
	if len(s) != 19 {
		t.Errorf("len(s) = %d, want 19 (bytes)", len(s))
	}
	if utf8.RuneCountInString(s) != 10 {
		t.Errorf("RuneCount = %d, want 10 (runes)", utf8.RuneCountInString(s))
	}
	t.Logf("bytes=%d, runes=%d — вот почему наивный цикл по байтам ломается", len(s), utf8.RuneCountInString(s))
}

func TestToLower(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"ascii", "Hello, World!", "hello, world!"},
		{"already lower", "hello", "hello"},
		{"cyrillic", "Привет, МИР!", "привет, мир!"},
		{"mixed with emoji", "Привет, 🌍!", "привет, 🌍!"},
		{"digits and punctuation", "ABC123!?", "abc123!?"},
		{"capital sharp s", "ẞ", "ß"},
		{"invalid utf8", "\xff", "\uFFFD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLower(tt.in)
			if got != tt.want {
				t.Errorf("ToLower(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSplitToRunes(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []rune
	}{
		{"empty", "", []rune{}},
		{"ascii", "abc", []rune{'a', 'b', 'c'}},
		{"cyrillic", "Привет", []rune{'П', 'р', 'и', 'в', 'е', 'т'}},
		{"with emoji", "a🌍b", []rune{'a', '🌍', 'b'}},
		{"full phrase", "Привет, 🌍!", []rune{'П', 'р', 'и', 'в', 'е', 'т', ',', ' ', '🌍', '!'}},
		{"invalid utf8", "\xff", []rune{utf8.RuneError}},
		{"invalid utf8 mixed", "a\xffb", []rune{'a', utf8.RuneError, 'b'}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitToRunes(tt.in)
			if !slices.Equal(got, tt.want) {
				t.Errorf("SplitToRunes(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestCleanNonLetters(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"only letters", "HelloWorld", "HelloWorld"},
		{"only digits", "12345", ""},
		{"only punctuation", "!@#$%^&*()", ""},
		{"with spaces", "Hello World", "HelloWorld"},
		{"with punctuation", "Hello, World!", "HelloWorld"},
		{"cyrillic with punctuation", "Привет, мир!", "Приветмир"},
		{"emoji removed", "Привет, 🌍!", "Привет"},
		{"mixed", "abc123 你好, 世界!", "abc你好世界"},
		{"punctuation removed", "a_b-c", "abc"},
		{"digits cyrillic", "тест123тест", "тесттест"},
		{"invalid utf8", "\xff", ""},
		{"invalid utf8 mixed", "a\xffb", "ab"}, // U+FFFD отбрасывается, a и b остаются
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanNonLetters(tt.in)
			if got != tt.want {
				t.Errorf("CleanNonLetters(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// Бонус: проверяем, что CleanNonLetters действительно удаляет всё, что не IsLetter
func TestCleanNonLetters_AllResultsAreLetters(t *testing.T) {
	inputs := []string{
		"Привет, 🌍!",
		"Hello, World! 123",
		"a_b-c.d",
		"!!!",
		"",
		"日本語テスト！",
	}
	for _, in := range inputs {
		out := CleanNonLetters(in)
		for _, r := range out {
			if !unicode.IsLetter(r) {
				t.Errorf("CleanNonLetters(%q) = %q contains non-letter %q", in, out, r)
			}
		}
	}
}

func TestSplitWords(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "empty string",
			in:   "",
			want: []string{},
		},
		{
			name: "only spaces",
			in:   " ",
			want: []string{},
		},
		{
			name: "only tabs and newlines",
			in:   "\t\n",
			want: []string{},
		},
		{
			name: "single space between words",
			in:   "hello world",
			want: []string{"hello", "world"},
		},
		{
			name: "multiple spaces between words",
			in:   "hello    world",
			want: []string{"hello", "world"},
		},
		{
			name: "leading and trailing spaces",
			in:   "  hello world  ",
			want: []string{"hello", "world"},
		},
		{
			name: "tabs, newlines, mixed whitespace",
			in:   "\thello\tworld\nfoo\r\nbar ",
			want: []string{"hello", "world", "foo", "bar"},
		},
		{
			name: "cyrillic",
			in:   "привет мир",
			want: []string{"привет", "мир"},
		},
		{
			name: "emoji between words",
			in:   "hello 😀 world 🎉 test",
			want: []string{"hello", "😀", "world", "🎉", "test"},
		},
		{
			name: "non-breaking space U+00A0",
			// unicode.IsSpace('\u00A0') == true, поэтому strings.Fields режет по нему
			in:   "hello\u00A0world",
			want: []string{"hello", "world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitWords(tt.in)
			if len(got) != len(tt.want) {
				t.Fatalf("SplitWords(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("SplitWords(%q) = %#v, want %#v", tt.in, got, tt.want)
				}
			}
		})
	}
}

func TestSplitWords_EmptyNotEmptyNil(t *testing.T) {
	got := SplitWords("")
	if got == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %#v", got)
	}
}
