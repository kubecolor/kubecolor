package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_DescribePrinter_Print(t *testing.T) {
	tests := []struct {
		name         string
		tablePrinter *TablePrinter
		input        string
		expected     string
	}{
		{
			name:         "values can be colored by its type",
			tablePrinter: NewTablePrinter(true, testconfig.DarkTheme, nil),
			input: testutil.NewHereDoc(`
				Name:         nginx-lpv5x
				Namespace:    default
				Priority:     0
				Node:         minikube/172.17.0.3
				Ready:        true
				Start Time:   Sat, 10 Oct 2020 14:07:17 +0900
				Labels:       app=nginx
				Annotations:  <none>
				Containers:
				  container-1:
				    Environment Variables from:
				      anycm	ConfigMap  Optional: true
				      anysec	Secret     Optional: false
				Conditions:
				  Type              Status
				  Initialized       True
				  Ready             False
				  ContainersReady   True
				  PodScheduled      True
				Volumes:
				  kube-api-access-7fdrt:
				    ConfigMapOptional:       <nil>`),
			expected: testutil.NewHereDoc(`
				\e[33mName\e[0m:         \e[37mnginx-lpv5x\e[0m
				\e[33mNamespace\e[0m:    \e[37mdefault\e[0m
				\e[33mPriority\e[0m:     \e[35m0\e[0m
				\e[33mNode\e[0m:         \e[37mminikube/172.17.0.3\e[0m
				\e[33mReady\e[0m:        \e[32mtrue\e[0m
				\e[33mStart Time\e[0m:   \e[37mSat, 10 Oct 2020 14:07:17 +0900\e[0m
				\e[33mLabels\e[0m:       app=\e[37mnginx\e[0m
				\e[33mAnnotations\e[0m:  \e[33m<none>\e[0m
				\e[33mContainers\e[0m:
				  \e[37mcontainer-1\e[0m:
				    \e[33mEnvironment Variables from\e[0m:
				      anycm	ConfigMap  Optional: \e[32mtrue\e[0m
				      anysec	Secret     Optional: \e[31mfalse\e[0m
				\e[33mConditions\e[0m:
				  \e[37mType\e[0m              \e[37mStatus\e[0m
				  \e[37mInitialized\e[0m       \e[32mTrue\e[0m
				  \e[37mReady\e[0m             \e[31mFalse\e[0m
				  \e[37mContainersReady\e[0m   \e[32mTrue\e[0m
				  \e[37mPodScheduled\e[0m      \e[32mTrue\e[0m
				\e[33mVolumes\e[0m:
				  \e[37mkube-api-access-7fdrt\e[0m:
				    \e[33mConfigMapOptional\e[0m:       \e[33m<nil>\e[0m
			`),
		},
		{
			name:         "key color changes based on its indentation",
			tablePrinter: NewTablePrinter(true, testconfig.DarkTheme, nil),
			input: testutil.NewHereDoc(`
				IP:           172.18.0.7
				IPs:
				  IP:           172.18.0.7
				Controlled By:  ReplicaSet/nginx
				Containers:
				  nginx:
				    Container ID:   docker://2885230a30908c8a6bda5a5366619c730b25b994eea61c931bba08ef4a8c8593
				      Started:      Sat, 10 Oct 2020 14:07:44 +0900`),
			expected: testutil.NewHereDoc(`
				\e[33mIP\e[0m:           \e[37m172.18.0.7\e[0m
				\e[33mIPs\e[0m:
				  \e[37mIP\e[0m:           \e[37m172.18.0.7\e[0m
				\e[33mControlled By\e[0m:  \e[37mReplicaSet/nginx\e[0m
				\e[33mContainers\e[0m:
				  \e[37mnginx\e[0m:
				    \e[33mContainer ID\e[0m:   \e[37mdocker://2885230a30908c8a6bda5a5366619c730b25b994eea61c931bba08ef4a8c8593\e[0m
				      \e[37mStarted\e[0m:      \e[37mSat, 10 Oct 2020 14:07:44 +0900\e[0m
			`),
		},
		{
			name:         "table format in kubectl describe can be colored by describe",
			tablePrinter: NewTablePrinter(false, testconfig.DarkTheme, nil),
			input: testutil.NewHereDoc(`
				Conditions:
				  Type             Status  LastHeartbeatTime                 LastTransitionTime                Reason                       Message
				  ----             ------  -----------------                 ------------------                ------                       -------
				  MemoryPressure   False   Sun, 18 Oct 2020 12:00:54 +0900   Wed, 14 Oct 2020 09:28:18 +0900   KubeletHasSufficientMemory   kubelet has sufficient memory available
				  DiskPressure     False   Sun, 18 Oct 2020 12:00:54 +0900   Wed, 14 Oct 2020 09:28:18 +0900   KubeletHasNoDiskPressure     kubelet has no disk pressure
				Addresses:
				  InternalIP:  172.17.0.3
				  Hostname:    minikube
				Capacity:
				  cpu:                6
				  memory:             2036900Ki
				  pods:               110
				Allocatable:
				  cpu:                6
				  memory:             2036900Ki
				  pods:               110
				System Info:
				  Machine ID:                 55d2ccaefc9847c9a69356e7f3bd23f4
				  System UUID:                fe312784-2364-4bba-a55e-f56051539c21
				Non-terminated Pods:          (14 in total)
				  Namespace                   Name                                CPU Requests  CPU Limits  Memory Requests  Memory Limits  AGE
				  ---------                   ----                                ------------  ----------  ---------------  -------------  ---
				  default                     nginx-6799fc88d8-dnmv5              0 (0%)        0 (0%)      0 (0%)           0 (0%)         7d21h
				  default                     nginx-6799fc88d8-m8pbc              0 (0%)        0 (0%)      0 (0%)           0 (0%)         7d21h
				  default                     nginx-6799fc88d8-qdf9b              0 (0%)        0 (0%)      0 (0%)           0 (0%)         7d21h
				Allocated resources:
				  (Total limits may be over 100 percent, i.e., overcommitted.)
				  Resource           Requests    Limits
				  --------           --------    ------
				  cpu                650m (10%)  0 (0%)
				  memory             70Mi (3%)   170Mi (8%)
				Events:              <none>`),
			expected: testutil.NewHereDoc(`
				\e[33mConditions\e[0m:
				  \e[37mType\e[0m             \e[36mStatus\e[0m  \e[37mLastHeartbeatTime\e[0m                 \e[36mLastTransitionTime\e[0m                \e[37mReason\e[0m                       \e[36mMessage\e[0m
				  \e[37m----\e[0m             \e[36m------\e[0m  \e[37m-----------------\e[0m                 \e[36m------------------\e[0m                \e[37m------\e[0m                       \e[36m-------\e[0m
				  \e[37mMemoryPressure\e[0m   \e[36mFalse\e[0m   \e[37mSun, 18 Oct 2020 12:00:54 +0900\e[0m   \e[36mWed, 14 Oct 2020 09:28:18 +0900\e[0m   \e[37mKubeletHasSufficientMemory\e[0m   \e[36mkubelet has sufficient memory available\e[0m
				  \e[37mDiskPressure\e[0m     \e[36mFalse\e[0m   \e[37mSun, 18 Oct 2020 12:00:54 +0900\e[0m   \e[36mWed, 14 Oct 2020 09:28:18 +0900\e[0m   \e[37mKubeletHasNoDiskPressure\e[0m     \e[36mkubelet has no disk pressure\e[0m
				\e[33mAddresses\e[0m:
				  \e[37mInternalIP\e[0m:  \e[37m172.17.0.3\e[0m
				  \e[37mHostname\e[0m:    \e[37mminikube\e[0m
				\e[33mCapacity\e[0m:
				  \e[37mcpu\e[0m:                \e[35m6\e[0m
				  \e[37mmemory\e[0m:             \e[37m2036900Ki\e[0m
				  \e[37mpods\e[0m:               \e[35m110\e[0m
				\e[33mAllocatable\e[0m:
				  \e[37mcpu\e[0m:                \e[35m6\e[0m
				  \e[37mmemory\e[0m:             \e[37m2036900Ki\e[0m
				  \e[37mpods\e[0m:               \e[35m110\e[0m
				\e[33mSystem Info\e[0m:
				  \e[37mMachine ID\e[0m:                 \e[37m55d2ccaefc9847c9a69356e7f3bd23f4\e[0m
				  \e[37mSystem UUID\e[0m:                \e[37mfe312784-2364-4bba-a55e-f56051539c21\e[0m
				\e[33mNon-terminated Pods\e[0m:          \e[37m(14 in total)\e[0m
				  \e[37mNamespace\e[0m                   \e[36mName\e[0m                                \e[37mCPU Requests\e[0m  \e[36mCPU Limits\e[0m  \e[37mMemory Requests\e[0m  \e[36mMemory Limits\e[0m  \e[37mAGE\e[0m
				  \e[37m---------\e[0m                   \e[36m----\e[0m                                \e[37m------------\e[0m  \e[36m----------\e[0m  \e[37m---------------\e[0m  \e[36m-------------\e[0m  \e[37m---\e[0m
				  \e[37mdefault\e[0m                     \e[36mnginx-6799fc88d8-dnmv5\e[0m              \e[37m0 (0%)\e[0m        \e[36m0 (0%)\e[0m      \e[37m0 (0%)\e[0m           \e[36m0 (0%)\e[0m         \e[37m7d21h\e[0m
				  \e[37mdefault\e[0m                     \e[36mnginx-6799fc88d8-m8pbc\e[0m              \e[37m0 (0%)\e[0m        \e[36m0 (0%)\e[0m      \e[37m0 (0%)\e[0m           \e[36m0 (0%)\e[0m         \e[37m7d21h\e[0m
				  \e[37mdefault\e[0m                     \e[36mnginx-6799fc88d8-qdf9b\e[0m              \e[37m0 (0%)\e[0m        \e[36m0 (0%)\e[0m      \e[37m0 (0%)\e[0m           \e[36m0 (0%)\e[0m         \e[37m7d21h\e[0m
				\e[33mAllocated resources\e[0m:
				  \e[37m(Total limits may be over 100 percent, i.e., overcommitted.)\e[0m
				  \e[37mResource\e[0m           \e[36mRequests\e[0m    \e[37mLimits\e[0m
				  \e[37m--------\e[0m           \e[36m--------\e[0m    \e[37m------\e[0m
				  \e[37mcpu\e[0m                \e[36m650m (10%)\e[0m  \e[37m0 (0%)\e[0m
				  \e[37mmemory\e[0m             \e[36m70Mi (3%)\e[0m   \e[37m170Mi (8%)\e[0m
				\e[33mEvents\e[0m:              \e[33m<none>\e[0m
			`),
		},
		{
			name:         "table format in kubectl describe at the end",
			tablePrinter: NewTablePrinter(false, testconfig.DarkTheme, nil),
			input: testutil.NewHereDoc(`
				Name:         cert-manager:leaderelection
				Labels:       app=cert-manager
											app.kubernetes.io/version=v1.12.3
											some-label=false
				Annotations:  meta.helm.sh/release-name: cert-manager
											meta.helm.sh/release-namespace: nais-system
											some-annotation: true
				PolicyRule:
					Resources                   Non-Resource URLs  Resource Names             Verbs
					---------                   -----------------  --------------             -----
					leases.coordination.k8s.io  []                 []                         [create]
					leases.coordination.k8s.io  []                 [cert-manager-controller]  [get update patch]`),
			expected: testutil.NewHereDoc(`
				\e[33mName\e[0m:         \e[37mcert-manager:leaderelection\e[0m
				\e[33mLabels\e[0m:       app=\e[37mcert-manager\e[0m
											app.kubernetes.io/version=\e[37mv1.12.3\e[0m
											some-label=\e[31mfalse\e[0m
				\e[33mAnnotations\e[0m:  meta.helm.sh/release-name: \e[37mcert-manager\e[0m
											meta.helm.sh/release-namespace: \e[37mnais-system\e[0m
											some-annotation: \e[32mtrue\e[0m
				\e[33mPolicyRule\e[0m:
					\e[37mResources\e[0m                   \e[36mNon-Resource URLs\e[0m  \e[37mResource Names\e[0m             \e[36mVerbs\e[0m
					\e[37m---------\e[0m                   \e[36m-----------------\e[0m  \e[37m--------------\e[0m             \e[36m-----\e[0m
					\e[37mleases.coordination.k8s.io\e[0m  \e[36m[]\e[0m                 \e[37m[]\e[0m                         \e[36m[create]\e[0m
					\e[37mleases.coordination.k8s.io\e[0m  \e[36m[]\e[0m                 \e[37m[cert-manager-controller]\e[0m  \e[36m[get update patch]\e[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := DescribePrinter{TablePrinter: tt.tablePrinter}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
