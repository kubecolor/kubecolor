package printer

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_TablePrinter_Print(t *testing.T) {
	var (
		colorRed    = config.MustParseColor("red")
		colorYellow = config.MustParseColor("yellow")
	)
	tests := []struct {
		name           string
		colorDeciderFn func(index int, column string) (config.Color, bool)
		withHeader     bool
		theme          *config.Theme
		input          string
		expected       string
	}{
		{
			name:           "header is not colored - dark",
			colorDeciderFn: nil,
			withHeader:     true,
			theme:          testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS    RESTARTS   AGE
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				\e[37mNAME          READY   STATUS    RESTARTS   AGE\e[0m
				\e[37mnginx-dnmv5\e[0m   \e[36m1/1\e[0m     \e[37mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-m8pbc\e[0m   \e[36m1/1\e[0m     \e[37mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-qdf9b\e[0m   \e[36m1/1\e[0m     \e[37mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
			`),
		},
		{
			name:           "multiple headers",
			colorDeciderFn: nil,
			withHeader:     true,
			theme:          testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				NAME                         READY   STATUS    RESTARTS   AGE
				pod/nginx-8spn9              1/1     Running   1          19d
				pod/nginx-dplns              1/1     Running   1          19d
				pod/nginx-lpv5x              1/1     Running   1          19d

				NAME                               DESIRED   CURRENT   READY   AGE
				replicaset.apps/nginx              3         3         3       19d
				replicaset.apps/nginx-6799fc88d8   3         3         3       19d
			`),
			expected: testutil.NewHereDoc(`
				\e[37mNAME                         READY   STATUS    RESTARTS   AGE\e[0m
				\e[37mpod/nginx-8spn9\e[0m              \e[36m1/1\e[0m     \e[37mRunning\e[0m   \e[36m1\e[0m          \e[37m19d\e[0m
				\e[37mpod/nginx-dplns\e[0m              \e[36m1/1\e[0m     \e[37mRunning\e[0m   \e[36m1\e[0m          \e[37m19d\e[0m
				\e[37mpod/nginx-lpv5x\e[0m              \e[36m1/1\e[0m     \e[37mRunning\e[0m   \e[36m1\e[0m          \e[37m19d\e[0m

				\e[37mNAME                               DESIRED   CURRENT   READY   AGE\e[0m
				\e[37mreplicaset.apps/nginx\e[0m              \e[36m3\e[0m         \e[37m3\e[0m         \e[36m3\e[0m       \e[37m19d\e[0m
				\e[37mreplicaset.apps/nginx-6799fc88d8\e[0m   \e[36m3\e[0m         \e[37m3\e[0m         \e[36m3\e[0m       \e[37m19d\e[0m
			`),
		},
		{
			name:           "withheader=false, 1st line is not colored in header color but colored as a content of table",
			colorDeciderFn: nil,
			withHeader:     false,
			theme:          testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				\e[37mnginx-dnmv5\e[0m   \e[36m1/1\e[0m     \e[37mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-m8pbc\e[0m   \e[36m1/1\e[0m     \e[37mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-qdf9b\e[0m   \e[36m1/1\e[0m     \e[37mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				`),
		},
		{
			name:           "when darkBackground=false, color preset for light is used",
			colorDeciderFn: nil,
			withHeader:     true,
			theme:          testconfig.LightTheme,
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS    RESTARTS   AGE
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				\e[30mNAME          READY   STATUS    RESTARTS   AGE\e[0m
				\e[30mnginx-dnmv5\e[0m   \e[34m1/1\e[0m     \e[30mRunning\e[0m   \e[34m0\e[0m          \e[30m6d6h\e[0m
				\e[30mnginx-m8pbc\e[0m   \e[34m1/1\e[0m     \e[30mRunning\e[0m   \e[34m0\e[0m          \e[30m6d6h\e[0m
				\e[30mnginx-qdf9b\e[0m   \e[34m1/1\e[0m     \e[30mRunning\e[0m   \e[34m0\e[0m          \e[30m6d6h\e[0m
			`),
		},
		{
			name: "colorDeciderFn works",
			colorDeciderFn: func(_ int, column string) (config.Color, bool) {
				if column == "CrashLoopBackOff" {
					return colorRed, true
				}

				// When Readiness is "n/m" then yellow
				if strings.Count(column, "/") == 1 {
					if arr := strings.Split(column, "/"); arr[0] != arr[1] {
						_, e1 := strconv.Atoi(arr[0])
						_, e2 := strconv.Atoi(arr[1])
						if e1 == nil && e2 == nil { // check both is number
							return colorYellow, true
						}
					}

				}

				return config.Color{}, false
			},
			withHeader: true,
			theme:      testconfig.DarkTheme,
			// "CrashLoopBackOff" will be red, "0/1" will be yellow
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS             RESTARTS   AGE
				nginx-dnmv5   1/1     CrashLoopBackOff   0          6d6h
				nginx-m8pbc   1/1     Running            0          6d6h
				nginx-qdf9b   0/1     Running            0          6d6h`),
			expected: testutil.NewHereDoc(`
				\e[37mNAME          READY   STATUS             RESTARTS   AGE\e[0m
				\e[37mnginx-dnmv5\e[0m   \e[36m1/1\e[0m     \e[31mCrashLoopBackOff\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-m8pbc\e[0m   \e[36m1/1\e[0m     \e[37mRunning\e[0m            \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-qdf9b\e[0m   \e[33m0/1\e[0m     \e[37mRunning\e[0m            \e[36m0\e[0m          \e[37m6d6h\e[0m
			`),
		},
		{
			name:           "a table whose some parts are missing can be handled",
			colorDeciderFn: nil,
			withHeader:     true,
			theme:          testconfig.DarkTheme,
			input: testutil.NewHereDoc(`
				NAME                              SHORTNAMES   APIGROUP                       NAMESPACED   KIND
				bindings                                                                      true         Binding
				componentstatuses                 cs                                          false        ComponentStatus
				pods                              po                                          true         Pod
				podtemplates                                                                  true         PodTemplate
				replicationcontrollers            rc                                          true         ReplicationController
				resourcequotas                    quota                                       true         ResourceQuota
				secrets                                                                       true         Secret
				serviceaccounts                   sa                                          true         ServiceAccount
				services                          svc                                         true         Service
				mutatingwebhookconfigurations                  admissionregistration.k8s.io   false        MutatingWebhookConfiguration
				customresourcedefinitions         crd,crds     apiextensions.k8s.io           false        CustomResourceDefinition
				controllerrevisions                            apps                           true         ControllerRevision
				daemonsets                        ds           apps                           true         DaemonSet
				statefulsets                      sts          apps                           true         StatefulSet
				tokenreviews                                   authentication.k8s.io          false        TokenReview
			`),
			expected: testutil.NewHereDoc(`
				\e[37mNAME                              SHORTNAMES   APIGROUP                       NAMESPACED   KIND\e[0m
				\e[37mbindings\e[0m                                                                      \e[36mtrue\e[0m         \e[37mBinding\e[0m
				\e[37mcomponentstatuses\e[0m                 \e[36mcs\e[0m                                          \e[36mfalse\e[0m        \e[37mComponentStatus\e[0m
				\e[37mpods\e[0m                              \e[36mpo\e[0m                                          \e[36mtrue\e[0m         \e[37mPod\e[0m
				\e[37mpodtemplates\e[0m                                                                  \e[36mtrue\e[0m         \e[37mPodTemplate\e[0m
				\e[37mreplicationcontrollers\e[0m            \e[36mrc\e[0m                                          \e[36mtrue\e[0m         \e[37mReplicationController\e[0m
				\e[37mresourcequotas\e[0m                    \e[36mquota\e[0m                                       \e[36mtrue\e[0m         \e[37mResourceQuota\e[0m
				\e[37msecrets\e[0m                                                                       \e[36mtrue\e[0m         \e[37mSecret\e[0m
				\e[37mserviceaccounts\e[0m                   \e[36msa\e[0m                                          \e[36mtrue\e[0m         \e[37mServiceAccount\e[0m
				\e[37mservices\e[0m                          \e[36msvc\e[0m                                         \e[36mtrue\e[0m         \e[37mService\e[0m
				\e[37mmutatingwebhookconfigurations\e[0m                  \e[37madmissionregistration.k8s.io\e[0m   \e[36mfalse\e[0m        \e[37mMutatingWebhookConfiguration\e[0m
				\e[37mcustomresourcedefinitions\e[0m         \e[36mcrd,crds\e[0m     \e[37mapiextensions.k8s.io\e[0m           \e[36mfalse\e[0m        \e[37mCustomResourceDefinition\e[0m
				\e[37mcontrollerrevisions\e[0m                            \e[37mapps\e[0m                           \e[36mtrue\e[0m         \e[37mControllerRevision\e[0m
				\e[37mdaemonsets\e[0m                        \e[36mds\e[0m           \e[37mapps\e[0m                           \e[36mtrue\e[0m         \e[37mDaemonSet\e[0m
				\e[37mstatefulsets\e[0m                      \e[36msts\e[0m          \e[37mapps\e[0m                           \e[36mtrue\e[0m         \e[37mStatefulSet\e[0m
				\e[37mtokenreviews\e[0m                                   \e[37mauthentication.k8s.io\e[0m          \e[36mfalse\e[0m        \e[37mTokenReview\e[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := NewTablePrinter(tt.withHeader, tt.theme, tt.colorDeciderFn)
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
