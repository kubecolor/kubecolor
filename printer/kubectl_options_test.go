package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func Test_OptionsPrinter_Print(t *testing.T) {
	tests := []struct {
		name           string
		darkBackground bool
		input          string
		expected       string
	}{
		{
			name:           "successful",
			darkBackground: true,
			input: testutil.NewHereDoc(`
				The following options can be passed to any command:
				
				      --add-dir-header=false: If true, adds the file directory to the header of the log messages
				      --alsologtostderr=false: log to standard error as well as files
				      --as='': Username to impersonate for the operation
				      --as-group=[]: Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
				      --cache-dir='/home/dtyler/.kube/cache': Default cache directory
				      --certificate-authority='': Path to a cert file for the certificate authority
				      --client-certificate='': Path to a client certificate file for TLS
				      --client-key='': Path to a client key file for TLS
				      --cluster='': The name of the kubeconfig cluster to use
				      --context='': The name of the kubeconfig context to use
				      --insecure-skip-tls-verify=false: If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
				`),
			expected: testutil.NewHereDoc(`
				[36mThe following options can be passed to any command:[0m
				
				      [33m--add-dir-header=false[0m: [36mIf true, adds the file directory to the header of the log messages[0m
				      [33m--alsologtostderr=false[0m: [36mlog to standard error as well as files[0m
				      [33m--as=''[0m: [36mUsername to impersonate for the operation[0m
				      [33m--as-group=[][0m: [36mGroup to impersonate for the operation, this flag can be repeated to specify multiple groups.[0m
				      [33m--cache-dir='/home/dtyler/.kube/cache'[0m: [36mDefault cache directory[0m
				      [33m--certificate-authority=''[0m: [36mPath to a cert file for the certificate authority[0m
				      [33m--client-certificate=''[0m: [36mPath to a client certificate file for TLS[0m
				      [33m--client-key=''[0m: [36mPath to a client key file for TLS[0m
				      [33m--cluster=''[0m: [36mThe name of the kubeconfig cluster to use[0m
				      [33m--context=''[0m: [36mThe name of the kubeconfig context to use[0m
				      [33m--insecure-skip-tls-verify=false[0m: [36mIf true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := OptionsPrinter{DarkBackground: tt.darkBackground}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
