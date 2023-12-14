package printer

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_TablePrinter_Print(t *testing.T) {
	tests := []struct {
		name           string
		colorDeciderFn func(index int, column string) (color.Color, bool)
		withHeader     bool
		darkBackground bool
		input          string
		expected       string
	}{
		{
			name:           "header is not colored - dark",
			colorDeciderFn: nil,
			withHeader:     true,
			darkBackground: true,
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS    RESTARTS   AGE
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				[37mNAME          READY   STATUS    RESTARTS   AGE[0m
				[37mnginx-dnmv5[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-m8pbc[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-qdf9b[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
			`),
		},
		{
			name:           "multiple headers",
			colorDeciderFn: nil,
			withHeader:     true,
			darkBackground: true,
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
				[37mNAME                         READY   STATUS    RESTARTS   AGE[0m
				[37mpod/nginx-8spn9[0m              [36m1/1[0m     [37mRunning[0m   [36m1[0m          [37m19d[0m
				[37mpod/nginx-dplns[0m              [36m1/1[0m     [37mRunning[0m   [36m1[0m          [37m19d[0m
				[37mpod/nginx-lpv5x[0m              [36m1/1[0m     [37mRunning[0m   [36m1[0m          [37m19d[0m
				[37m[0m
				[37mNAME                               DESIRED   CURRENT   READY   AGE[0m
				[37mreplicaset.apps/nginx[0m              [36m3[0m         [37m3[0m         [36m3[0m       [37m19d[0m
				[37mreplicaset.apps/nginx-6799fc88d8[0m   [36m3[0m         [37m3[0m         [36m3[0m       [37m19d[0m
			`),
		},
		{
			name:           "withheader=false, 1st line is not colored in header color but colored as a content of table",
			colorDeciderFn: nil,
			withHeader:     false,
			darkBackground: true,
			input: testutil.NewHereDoc(`
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				[37mnginx-dnmv5[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-m8pbc[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-qdf9b[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
				`),
		},
		{
			name:           "when darkBackground=false, color preset for light is used",
			colorDeciderFn: nil,
			withHeader:     true,
			darkBackground: false,
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS    RESTARTS   AGE
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				[30mNAME          READY   STATUS    RESTARTS   AGE[0m
				[30mnginx-dnmv5[0m   [34m1/1[0m     [30mRunning[0m   [34m0[0m          [30m6d6h[0m
				[30mnginx-m8pbc[0m   [34m1/1[0m     [30mRunning[0m   [34m0[0m          [30m6d6h[0m
				[30mnginx-qdf9b[0m   [34m1/1[0m     [30mRunning[0m   [34m0[0m          [30m6d6h[0m
			`),
		},
		{
			name: "colorDeciderFn works",
			colorDeciderFn: func(_ int, column string) (color.Color, bool) {
				if column == "CrashLoopBackOff" {
					return color.Red, true
				}

				// When Readiness is "n/m" then yellow
				if strings.Count(column, "/") == 1 {
					if arr := strings.Split(column, "/"); arr[0] != arr[1] {
						_, e1 := strconv.Atoi(arr[0])
						_, e2 := strconv.Atoi(arr[1])
						if e1 == nil && e2 == nil { // check both is number
							return color.Yellow, true
						}
					}

				}

				return 0, false
			},
			withHeader:     true,
			darkBackground: true,
			// "CrashLoopBackOff" will be red, "0/1" will be yellow
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS             RESTARTS   AGE
				nginx-dnmv5   1/1     CrashLoopBackOff   0          6d6h
				nginx-m8pbc   1/1     Running            0          6d6h
				nginx-qdf9b   0/1     Running            0          6d6h`),
			expected: testutil.NewHereDoc(`
				[37mNAME          READY   STATUS             RESTARTS   AGE[0m
				[37mnginx-dnmv5[0m   [36m1/1[0m     [31mCrashLoopBackOff[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-m8pbc[0m   [36m1/1[0m     [37mRunning[0m            [36m0[0m          [37m6d6h[0m
				[37mnginx-qdf9b[0m   [33m0/1[0m     [37mRunning[0m            [36m0[0m          [37m6d6h[0m
			`),
		},
		{
			name:           "a table whose some parts are missing can be handled",
			colorDeciderFn: nil,
			withHeader:     true,
			darkBackground: true,
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
				[37mNAME                              SHORTNAMES   APIGROUP                       NAMESPACED   KIND[0m
				[37mbindings[0m                                                                      [36mtrue[0m         [37mBinding[0m
				[37mcomponentstatuses[0m                 [36mcs[0m                                          [36mfalse[0m        [37mComponentStatus[0m
				[37mpods[0m                              [36mpo[0m                                          [36mtrue[0m         [37mPod[0m
				[37mpodtemplates[0m                                                                  [36mtrue[0m         [37mPodTemplate[0m
				[37mreplicationcontrollers[0m            [36mrc[0m                                          [36mtrue[0m         [37mReplicationController[0m
				[37mresourcequotas[0m                    [36mquota[0m                                       [36mtrue[0m         [37mResourceQuota[0m
				[37msecrets[0m                                                                       [36mtrue[0m         [37mSecret[0m
				[37mserviceaccounts[0m                   [36msa[0m                                          [36mtrue[0m         [37mServiceAccount[0m
				[37mservices[0m                          [36msvc[0m                                         [36mtrue[0m         [37mService[0m
				[37mmutatingwebhookconfigurations[0m                  [37madmissionregistration.k8s.io[0m   [36mfalse[0m        [37mMutatingWebhookConfiguration[0m
				[37mcustomresourcedefinitions[0m         [36mcrd,crds[0m     [37mapiextensions.k8s.io[0m           [36mfalse[0m        [37mCustomResourceDefinition[0m
				[37mcontrollerrevisions[0m                            [37mapps[0m                           [36mtrue[0m         [37mControllerRevision[0m
				[37mdaemonsets[0m                        [36mds[0m           [37mapps[0m                           [36mtrue[0m         [37mDaemonSet[0m
				[37mstatefulsets[0m                      [36msts[0m          [37mapps[0m                           [36mtrue[0m         [37mStatefulSet[0m
				[37mtokenreviews[0m                                   [37mauthentication.k8s.io[0m          [36mfalse[0m        [37mTokenReview[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := NewTablePrinter(tt.withHeader, tt.darkBackground, tt.colorDeciderFn)
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
