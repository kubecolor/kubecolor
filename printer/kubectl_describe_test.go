package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/color"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_DescribePrinter_Print(t *testing.T) {
	tests := []struct {
		name           string
		darkBackground bool
		tablePrinter   *TablePrinter
		input          string
		expected       string
	}{
		{
			name:           "values can be colored by its type",
			darkBackground: true,
			tablePrinter:   NewTablePrinter(true, true, color.NewTheme(color.PresetDark), nil),
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
				[33mName[0m:         [37mnginx-lpv5x[0m
				[33mNamespace[0m:    [37mdefault[0m
				[33mPriority[0m:     [35m0[0m
				[33mNode[0m:         [37mminikube/172.17.0.3[0m
				[33mReady[0m:        [32mtrue[0m
				[33mStart Time[0m:   [37mSat, 10 Oct 2020 14:07:17 +0900[0m
				[33mLabels[0m:       app=[37mnginx[0m
				[33mAnnotations[0m:  [33m<none>[0m
				[33mContainers[0m:
				  [37mcontainer-1[0m:
				    [33mEnvironment Variables from[0m:
				      anycm	ConfigMap  Optional: [32mtrue[0m
				      anysec	Secret     Optional: [31mfalse[0m
				[33mConditions[0m:
				  [37mType[0m              [37mStatus[0m
				  [37mInitialized[0m       [32mTrue[0m
				  [37mReady[0m             [31mFalse[0m
				  [37mContainersReady[0m   [32mTrue[0m
				  [37mPodScheduled[0m      [32mTrue[0m
				[33mVolumes[0m:
				  [37mkube-api-access-7fdrt[0m:
				    [33mConfigMapOptional[0m:       [33m<nil>[0m
			`),
		},
		{
			name:           "key color changes based on its indentation",
			darkBackground: true,
			tablePrinter:   NewTablePrinter(true, true, color.NewTheme(color.PresetDark), nil),
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
				[33mIP[0m:           [37m172.18.0.7[0m
				[33mIPs[0m:
				  [37mIP[0m:           [37m172.18.0.7[0m
				[33mControlled By[0m:  [37mReplicaSet/nginx[0m
				[33mContainers[0m:
				  [37mnginx[0m:
				    [33mContainer ID[0m:   [37mdocker://2885230a30908c8a6bda5a5366619c730b25b994eea61c931bba08ef4a8c8593[0m
				      [37mStarted[0m:      [37mSat, 10 Oct 2020 14:07:44 +0900[0m
			`),
		},
		{
			name:           "table format in kubectl describe can be colored by describe",
			darkBackground: true,
			tablePrinter:   NewTablePrinter(false, true, color.NewTheme(color.PresetDark), nil),
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
				[33mConditions[0m:
				  [37mType[0m             [36mStatus[0m  [37mLastHeartbeatTime[0m                 [36mLastTransitionTime[0m                [37mReason[0m                       [36mMessage[0m
				  [37m----[0m             [36m------[0m  [37m-----------------[0m                 [36m------------------[0m                [37m------[0m                       [36m-------[0m
				  [37mMemoryPressure[0m   [36mFalse[0m   [37mSun, 18 Oct 2020 12:00:54 +0900[0m   [36mWed, 14 Oct 2020 09:28:18 +0900[0m   [37mKubeletHasSufficientMemory[0m   [36mkubelet has sufficient memory available[0m
				  [37mDiskPressure[0m     [36mFalse[0m   [37mSun, 18 Oct 2020 12:00:54 +0900[0m   [36mWed, 14 Oct 2020 09:28:18 +0900[0m   [37mKubeletHasNoDiskPressure[0m     [36mkubelet has no disk pressure[0m
				[33mAddresses[0m:
				  [37mInternalIP[0m:  [37m172.17.0.3[0m
				  [37mHostname[0m:    [37mminikube[0m
				[33mCapacity[0m:
				  [37mcpu[0m:                [35m6[0m
				  [37mmemory[0m:             [37m2036900Ki[0m
				  [37mpods[0m:               [35m110[0m
				[33mAllocatable[0m:
				  [37mcpu[0m:                [35m6[0m
				  [37mmemory[0m:             [37m2036900Ki[0m
				  [37mpods[0m:               [35m110[0m
				[33mSystem Info[0m:
				  [37mMachine ID[0m:                 [37m55d2ccaefc9847c9a69356e7f3bd23f4[0m
				  [37mSystem UUID[0m:                [37mfe312784-2364-4bba-a55e-f56051539c21[0m
				[33mNon-terminated Pods[0m:          [37m(14 in total)[0m
				  [37mNamespace[0m                   [36mName[0m                                [37mCPU Requests[0m  [36mCPU Limits[0m  [37mMemory Requests[0m  [36mMemory Limits[0m  [37mAGE[0m
				  [37m---------[0m                   [36m----[0m                                [37m------------[0m  [36m----------[0m  [37m---------------[0m  [36m-------------[0m  [37m---[0m
				  [37mdefault[0m                     [36mnginx-6799fc88d8-dnmv5[0m              [37m0 (0%)[0m        [36m0 (0%)[0m      [37m0 (0%)[0m           [36m0 (0%)[0m         [37m7d21h[0m
				  [37mdefault[0m                     [36mnginx-6799fc88d8-m8pbc[0m              [37m0 (0%)[0m        [36m0 (0%)[0m      [37m0 (0%)[0m           [36m0 (0%)[0m         [37m7d21h[0m
				  [37mdefault[0m                     [36mnginx-6799fc88d8-qdf9b[0m              [37m0 (0%)[0m        [36m0 (0%)[0m      [37m0 (0%)[0m           [36m0 (0%)[0m         [37m7d21h[0m
				[33mAllocated resources[0m:
				  [37m(Total limits may be over 100 percent, i.e., overcommitted.)[0m
				  [37mResource[0m           [36mRequests[0m    [37mLimits[0m
				  [37m--------[0m           [36m--------[0m    [37m------[0m
				  [37mcpu[0m                [36m650m (10%)[0m  [37m0 (0%)[0m
				  [37mmemory[0m             [36m70Mi (3%)[0m   [37m170Mi (8%)[0m
				[33mEvents[0m:              [33m<none>[0m
			`),
		},
		{
			name:           "table format in kubectl describe at the end",
			darkBackground: true,
			tablePrinter:   NewTablePrinter(false, true, nil),
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
				[33mName[0m:         [37mcert-manager:leaderelection[0m
				[33mLabels[0m:       app=[37mcert-manager[0m
											app.kubernetes.io/version=[37mv1.12.3[0m
											some-label=[31mfalse[0m
				[33mAnnotations[0m:  meta.helm.sh/release-name: [37mcert-manager[0m
											meta.helm.sh/release-namespace: [37mnais-system[0m
											some-annotation: [32mtrue[0m
				[33mPolicyRule[0m:
					[37mResources[0m                   [36mNon-Resource URLs[0m  [37mResource Names[0m             [36mVerbs[0m
					[37m---------[0m                   [36m-----------------[0m  [37m--------------[0m             [36m-----[0m
					[37mleases.coordination.k8s.io[0m  [36m[][0m                 [37m[][0m                         [36m[create][0m
					[37mleases.coordination.k8s.io[0m  [36m[][0m                 [37m[cert-manager-controller][0m  [36m[get update patch][0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := DescribePrinter{DarkBackground: tt.darkBackground, TablePrinter: tt.tablePrinter}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
