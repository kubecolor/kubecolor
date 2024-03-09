package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_VersionPrinter_Print(t *testing.T) {
	tests := []struct {
		name        string
		themePreset config.Preset
		recursive   bool
		input       string
		expected    string
	}{
		{
			name:        "go struct dump can be colorized",
			themePreset: config.PresetDark,
			input: testutil.NewHereDoc(`
				Client Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.3", GitCommit:"1e11e4a2108024935ecfcb2912226cedeafd99df", GitTreeState:"clean", BuildDate:"2020-10-14T18:49:28Z", GoVersion:"go1.15.2", Compiler:"gc", Platform:"darwin/amd64"}
				Server Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.2", GitCommit:"f5743093fd1c663cb0cbc89748f730662345d44d", GitTreeState:"clean", BuildDate:"2020-09-16T13:32:58Z", GoVersion:"go1.15", Compiler:"gc", Platform:"linux/amd64"}`),
			expected: testutil.NewHereDoc(`
				[33mClient Version[0m: [37mversion.Info[0m{[33mMajor[0m:"[37m1[0m", [33mMinor[0m:"[37m19[0m", [33mGitVersion[0m:"[37mv1.19.3[0m", [33mGitCommit[0m:"[37m1e11e4a2108024935ecfcb2912226cedeafd99df[0m", [33mGitTreeState[0m:"[37mclean[0m", [33mBuildDate[0m:"[37m2020-10-14T18:49:28Z[0m", [33mGoVersion[0m:"[37mgo1.15.2[0m", [33mCompiler[0m:"[37mgc[0m", [33mPlatform[0m:"[37mdarwin/amd64[0m"}
				[33mServer Version[0m: [37mversion.Info[0m{[33mMajor[0m:"[37m1[0m", [33mMinor[0m:"[37m19[0m", [33mGitVersion[0m:"[37mv1.19.2[0m", [33mGitCommit[0m:"[37mf5743093fd1c663cb0cbc89748f730662345d44d[0m", [33mGitTreeState[0m:"[37mclean[0m", [33mBuildDate[0m:"[37m2020-09-16T13:32:58Z[0m", [33mGoVersion[0m:"[37mgo1.15[0m", [33mCompiler[0m:"[37mgc[0m", [33mPlatform[0m:"[37mlinux/amd64[0m"}
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := VersionPrinter{
				Theme: config.NewTheme(tt.themePreset),
			}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}

func Test_VersionClientPrinter_Print(t *testing.T) {
	tests := []struct {
		name        string
		themePreset config.Preset
		input       string
		expected    string
	}{
		{
			name:        "--client can be colorized",
			themePreset: config.PresetDark,
			input: testutil.NewHereDoc(`
				Client Version: v1.19.3
				Server Version: v1.19.2`),
			expected: testutil.NewHereDoc(`
				[33mClient Version[0m: [37mv1.19.3[0m
				[33mServer Version[0m: [37mv1.19.2[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := VersionClientPrinter{
				Theme: config.NewTheme(tt.themePreset),
			}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
