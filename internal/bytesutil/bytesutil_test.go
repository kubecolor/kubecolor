package bytesutil

import "testing"

func TestCountColumns(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "empty",
			input: "",
			want:  0,
		},

		{
			name:  "only space",
			input: "       ",
			want:  0,
		},

		{
			name:  "single",
			input: "foo",
			want:  1,
		},

		{
			name:  "three/narrow spacing",
			input: "foo  bar  moo",
			want:  3,
		},

		{
			// This is where an implementation using [strings.Count] would fail.
			name:  "three/excessive spacing",
			input: "foo                bar             moo",
			want:  3,
		},
	}

	const spaceCharset = " \t"

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CountColumns([]byte(tc.input), spaceCharset)
			if got != tc.want {
				t.Errorf("Want %d, got %d\nInput: %q", tc.want, got, tc.input)
			}
		})
	}
}
