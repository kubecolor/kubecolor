package logscan

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestScanner_inputOutputMatches(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty",
			input: "",
		},
		{
			name:  "empty lines",
			input: "\n\n\n",
		},
		{
			name:  "single word",
			input: "helloworld\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			scanner := NewScanner(strings.NewReader(tc.input))

			var buf bytes.Buffer
			for scanner.Scan() {
				buf.WriteString(scanner.Token().Text)
			}

			testutil.MustEqual(t, tc.input, buf.String())
		})
	}
}

func TestScanner_tokens(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "empty",
			input: "",
			want:  nil,
		},
		{
			name:  "pre-formatted",
			input: "\033e[33m[NOTE]\033[0m this line already has colored output\n",
			want: []Token{
				{Kind: KindPreformatted, Text: "\033e[33m[NOTE]\033[0m this line already has colored output"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "single line",
			input: "\n",
			want: []Token{
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "key=value",
			input: "key=value other\n",
			want: []Token{
				{Kind: KindKey, Text: "key"},
				{Kind: KindUnknown, Text: "="},
				{Kind: KindValue, Text: "value"},
				{Kind: KindUnknown, Text: " "},
				{Kind: KindUnknown, Text: "other"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "key=value in parenthases",
			input: "(key=value)\n",
			want: []Token{
				{Kind: KindParenthases, Text: "("},
				{Kind: KindKey, Text: "key"},
				{Kind: KindUnknown, Text: "="},
				{Kind: KindValue, Text: "value"},
				{Kind: KindParenthases, Text: ")"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "key=value with value in parenthases",
			input: "key=(value with spaces)\n",
			want: []Token{
				{Kind: KindKey, Text: "key"},
				{Kind: KindUnknown, Text: "="},
				{Kind: KindValue, Text: "(value with spaces)"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "key=value with value in quotes",
			input: "key=\"value with spaces\"\n",
			want: []Token{
				{Kind: KindKey, Text: "key"},
				{Kind: KindUnknown, Text: "="},
				{Kind: KindValue, Text: "\"value with spaces\""},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "date UTC",
			input: "2024-08-03T12:38:44.049832713Z\n",
			want: []Token{
				{Kind: KindDate, Text: "2024-08-03T12:38:44.049832713Z"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "date timezone",
			input: "2024-08-03T12:38:44.049832713+02:00\n",
			want: []Token{
				{Kind: KindDate, Text: "2024-08-03T12:38:44.049832713+02:00"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "guid dashes",
			input: "70d5707e-b07b-41c3-9411-cad84c6db764\n",
			want: []Token{
				{Kind: KindGUID, Text: "70d5707e-b07b-41c3-9411-cad84c6db764"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "guid no dashes",
			input: "70d5707eb07b41c39411cad84c6db764\n",
			want: []Token{
				{Kind: KindGUID, Text: "70d5707eb07b41c39411cad84c6db764"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			scanner := NewScanner(strings.NewReader(tc.input))

			var tokens []Token
			for scanner.Scan() {
				tokens = append(tokens, scanner.Token())
			}

			testutil.MustEqual(t, tc.want, tokens)
		})
	}
}

func TestReadWord(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "letters",
			input: "hello world",
			want:  "hello",
		},
		{
			name:  "spaces",
			input: "\t hello world",
			want:  "\t ",
		},
		{
			name:  "punctuation",
			input: "(hello-world!)",
			want:  "(hello-world!)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			word := readWord([]byte(tc.input))
			testutil.MustEqual(t, tc.want, string(word), "input: "+tc.input)
		})
	}
}

func TestReadParenthases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "empty group",
			input: "()",
			want:  "()",
		},
		{
			name:  "single",
			input: "(hello world)",
			want:  "(hello world)",
		},
		{
			name:  "one after the other",
			input: "(hello world) (another one here)",
			want:  "(hello world)",
		},
		{
			name:  "nested",
			input: "(hello (another one here) world)",
			want:  "(hello (another one here) world)",
		},
		{
			name:  "never closed",
			input: "(hello world",
			want:  "",
		},
		{
			name:  "never closed nested",
			input: "(hello world (another one here)",
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			word := readParenthases([]byte(tc.input), '(', ')')
			testutil.MustEqual(t, tc.want, string(word), "input: "+tc.input)
		})
	}
}

func TestReadQuoted(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "empty quote",
			input: `""`,
			want:  `""`,
		},
		{
			name:  "simple quote",
			input: `"hello world"`,
			want:  `"hello world"`,
		},
		{
			name:  "ignores after",
			input: `"hello world" foo bar`,
			want:  `"hello world"`,
		},
		{
			name:  "escaped quote",
			input: `"hello \" world"`,
			want:  `"hello \" world"`,
		},
		{
			name:  "escaped escape",
			input: `"hello \\" world"`,
			want:  `"hello \\"`,
		},
		{
			name:  "escaped quote with escaped escape",
			input: `"hello \\\" world"`,
			want:  `"hello \\\" world"`,
		},
		{
			name:  "single quotes",
			input: `'hello " world'`,
			want:  `'hello " world'`,
		},
		{
			name:  "tick quotes",
			input: "`hello \" world`",
			want:  "`hello \" world`",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			quoted := readQuoted([]byte(tc.input))
			testutil.MustEqual(t, tc.want, string(quoted), "input: "+tc.input)
		})
	}
}
