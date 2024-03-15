package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_VersionPrinter_Print(t *testing.T) {
	tests := []struct {
		name      string
		theme     *config.Theme
		recursive bool
		input     string
		expected  string
	}{
		{
			name:  "go struct dump can be colorized",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				Client Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.3", GitCommit:"1e11e4a2108024935ecfcb2912226cedeafd99df", GitTreeState:"clean", BuildDate:"2020-10-14T18:49:28Z", GoVersion:"go1.15.2", Compiler:"gc", Platform:"darwin/amd64"}
				Server Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.2", GitCommit:"f5743093fd1c663cb0cbc89748f730662345d44d", GitTreeState:"clean", BuildDate:"2020-09-16T13:32:58Z", GoVersion:"go1.15", Compiler:"gc", Platform:"linux/amd64"}`),
			expected: testutil.NewHereDoc(`
				\e[33mClient Version\e[0m: \e[37mversion.Info\e[0m{\e[33mMajor\e[0m:"\e[37m1\e[0m", \e[33mMinor\e[0m:"\e[37m19\e[0m", \e[33mGitVersion\e[0m:"\e[37mv1.19.3\e[0m", \e[33mGitCommit\e[0m:"\e[37m1e11e4a2108024935ecfcb2912226cedeafd99df\e[0m", \e[33mGitTreeState\e[0m:"\e[37mclean\e[0m", \e[33mBuildDate\e[0m:"\e[37m2020-10-14T18:49:28Z\e[0m", \e[33mGoVersion\e[0m:"\e[37mgo1.15.2\e[0m", \e[33mCompiler\e[0m:"\e[37mgc\e[0m", \e[33mPlatform\e[0m:"\e[37mdarwin/amd64\e[0m"}
				\e[33mServer Version\e[0m: \e[37mversion.Info\e[0m{\e[33mMajor\e[0m:"\e[37m1\e[0m", \e[33mMinor\e[0m:"\e[37m19\e[0m", \e[33mGitVersion\e[0m:"\e[37mv1.19.2\e[0m", \e[33mGitCommit\e[0m:"\e[37mf5743093fd1c663cb0cbc89748f730662345d44d\e[0m", \e[33mGitTreeState\e[0m:"\e[37mclean\e[0m", \e[33mBuildDate\e[0m:"\e[37m2020-09-16T13:32:58Z\e[0m", \e[33mGoVersion\e[0m:"\e[37mgo1.15\e[0m", \e[33mCompiler\e[0m:"\e[37mgc\e[0m", \e[33mPlatform\e[0m:"\e[37mlinux/amd64\e[0m"}
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
				Theme: tt.theme,
			}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}

func Test_VersionClientPrinter_Print(t *testing.T) {
	tests := []struct {
		name     string
		theme    *config.Theme
		input    string
		expected string
	}{
		{
			name:  "--client can be colorized",
			theme: testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				Client Version: v1.19.3
				Server Version: v1.19.2`),
			expected: testutil.NewHereDoc(`
				\e[33mClient Version\e[0m: \e[37mv1.19.3\e[0m
				\e[33mServer Version\e[0m: \e[37mv1.19.2\e[0m
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
				Theme: tt.theme,
			}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
