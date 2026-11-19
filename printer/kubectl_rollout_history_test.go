package printer

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_parseRolloutHistoryResource(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		wantName     string
		wantRevision string
		wantOK       bool
	}{
		{"name with trailing space", "deployment.apps/nginx ", "deployment.apps/nginx", "", true},
		{"name without trailing space", "deployment.apps/nginx", "deployment.apps/nginx", "", true},
		{"core group resource", "pod/nginx ", "pod/nginx", "", true},
		{"with revision", "deployment.apps/nginx with revision #2", "deployment.apps/nginx", "2", true},
		{"with multi-digit revision", "daemonset.apps/fluentd with revision #15", "daemonset.apps/fluentd", "15", true},

		{"empty line", "", "", "", false},
		{"indented line", "  Labels:       app=nginx", "", "", false},
		{"table header", "REVISION  CHANGE-CAUSE", "", "", false},
		{"table row", "2         kubectl set image deployment/nginx nginx=nginx:1.27", "", "", false},
		{"name followed by other text", "deployment.apps/nginx rolled back", "", "", false},
		{"name followed by tab", "deployment.apps/nginx\twith revision #2", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, revision, ok := parseRolloutHistoryResource(tt.line)
			testutil.Equal(t, tt.wantName, name, "name")
			testutil.Equal(t, tt.wantRevision, revision, "revision")
			testutil.Equal(t, tt.wantOK, ok, "ok")
		})
	}
}

func TestRolloutHistoryPrinter_keepsOutputWithoutColors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "table",
			// kubectl ends the resource line with a space when --revision is not set
			input: "deployment.apps/nginx \n" +
				"REVISION  CHANGE-CAUSE\n" +
				"1         <none>\n" +
				"2         <none>\n",
		},
		{
			name: "multiple resources",
			input: "deployment.apps/nginx \n" +
				"REVISION  CHANGE-CAUSE\n" +
				"1         <none>\n" +
				"\n" +
				"daemonset.apps/fluentd \n" +
				"REVISION  CHANGE-CAUSE\n" +
				"3         <none>\n",
		},
		{
			name: "with revision",
			input: testutil.NewHereDoc(`
				deployment.apps/nginx with revision #2
				Pod Template:
				  Labels:  app=nginx
				  Containers:
				   nginx:
				    Image:         nginx:1.27
				  Node-Selectors:  <none>
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			printer := RolloutHistoryPrinter{Theme: testconfig.NullTheme}
			printer.Print(strings.NewReader(tt.input), &out)
			testutil.Equal(t, tt.input, out.String())
		})
	}
}

func TestRolloutHistoryPrinter_fail(t *testing.T) {
	var logBuf bytes.Buffer
	testutil.SetTestLogger(t, &logBuf)

	var outBuf bytes.Buffer
	printer := RolloutHistoryPrinter{Theme: testconfig.NullTheme}
	printer.Print(testutil.DummyReader{ReadFunc: func(b []byte) (int, error) { return 0, errors.New("test") }}, &outBuf)

	testutil.Equal(t, "", outBuf.String(), "output")
	testutil.Equal(t, "level=ERROR msg=\"Failed to print rollout history output.\" error=test\n", logBuf.String(), "logs")
}
