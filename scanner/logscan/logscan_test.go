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
			name:  "parenthases with word",
			input: "(foo)\n",
			want: []Token{
				{Kind: KindParenthases, Text: "("},
				{Kind: KindUnknown, Text: "foo"},
				{Kind: KindParenthases, Text: ")"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			// https://github.com/kubecolor/kubecolor/issues/167
			name:  "parenthases with words",
			input: "(foo bar)\n",
			want: []Token{
				{Kind: KindParenthases, Text: "("},
				{Kind: KindUnknown, Text: "foo"},
				{Kind: KindUnknown, Text: " "},
				{Kind: KindUnknown, Text: "bar"},
				{Kind: KindParenthases, Text: ")"},
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

		// JSON
		{
			name:  "json empty object",
			input: "{}\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json empty object with spaces",
			input: "{  }\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindUnknown, Text: "  "},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json object with single key",
			input: `{"key":"value"}` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindValue, Text: `"value"`},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json object with single key with spaces",
			input: `{  "key"  :  "value"  }` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindUnknown, Text: "  "},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: "  :"},
				{Kind: KindUnknown, Text: "  "},
				{Kind: KindValue, Text: `"value"`},
				{Kind: KindUnknown, Text: "  "},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json object with two keys",
			input: `{"key":"value","key2":"value2"}` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindValue, Text: `"value"`},
				{Kind: KindUnknown, Text: ","},
				{Kind: KindKey, Text: `"key2"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindValue, Text: `"value2"`},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json object with boolean",
			input: `{"key":true}` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindValue, Text: `true`},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json object with number",
			input: `{"key":123.456}` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindValue, Text: `123.456`},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json empty array",
			input: `{"key":[]}` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindParenthases, Text: "["},
				{Kind: KindParenthases, Text: "]"},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json empty array with spaces",
			input: `{"key":[  ]}` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindParenthases, Text: "["},
				{Kind: KindUnknown, Text: "  "},
				{Kind: KindParenthases, Text: "]"},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json array with single item",
			input: `{"key":[1]}` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindParenthases, Text: "["},
				{Kind: KindValue, Text: "1"},
				{Kind: KindParenthases, Text: "]"},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json array with two items",
			input: `{"key":[1,2]}` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindParenthases, Text: "["},
				{Kind: KindValue, Text: "1"},
				{Kind: KindUnknown, Text: ","},
				{Kind: KindValue, Text: "2"},
				{Kind: KindParenthases, Text: "]"},
				{Kind: KindParenthases, Text: "}"},
				{Kind: KindNewline, Text: "\n"},
			},
		},
		{
			name:  "json array with two items with spaces",
			input: `{"key":[  1  ,  2  ]}` + "\n",
			want: []Token{
				{Kind: KindParenthases, Text: "{"},
				{Kind: KindKey, Text: `"key"`},
				{Kind: KindUnknown, Text: ":"},
				{Kind: KindParenthases, Text: "["},
				{Kind: KindUnknown, Text: "  "},
				{Kind: KindValue, Text: "1"},
				{Kind: KindUnknown, Text: "  "},
				{Kind: KindUnknown, Text: ","},
				{Kind: KindUnknown, Text: "  "},
				{Kind: KindValue, Text: "2"},
				{Kind: KindUnknown, Text: "  "},
				{Kind: KindParenthases, Text: "]"},
				{Kind: KindParenthases, Text: "}"},
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

func TestReadLetters(t *testing.T) {
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
			name:  "single word",
			input: `foo`,
			want:  `foo`,
		},
		{
			name:  "up until numbers",
			input: `foo123`,
			want:  `foo`,
		},
		{
			name:  "up until symbol",
			input: `foo}`,
			want:  `foo`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			quoted := readLetters([]byte(tc.input))
			testutil.MustEqual(t, tc.want, string(quoted), "input: "+tc.input)
		})
	}
}
