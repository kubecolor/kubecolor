package printer

import (
	"bytes"
	"strings"
	"testing"

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
			tablePrinter:   nil,
			input: testutil.NewHereDoc(`
				Name:         nginx-lpv5x
				Namespace:    default
				Priority:     0
				Node:         minikube/172.17.0.3
				Ready:        true
				Start Time:   Sat, 10 Oct 2020 14:07:17 +0900
				Labels:       app=nginx
				Annotations:  <none>`),
			expected: testutil.NewHereDoc(`
				[33mName[0m:         [36mnginx-lpv5x[0m
				[33mNamespace[0m:    [36mdefault[0m
				[33mPriority[0m:     [35m0[0m
				[33mNode[0m:         [36mminikube/172.17.0.3[0m
				[33mReady[0m:        [32mtrue[0m
				[33mStart Time[0m:   [36mSat, 10 Oct 2020 14:07:17 +0900[0m
				[33mLabels[0m:       [36mapp=nginx[0m
				[33mAnnotations[0m:  [33m<none>[0m
			`),
		},
		{
			name:           "key color changes based on its indentation",
			darkBackground: true,
			tablePrinter:   nil,
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
				[33mIP[0m:           [36m172.18.0.7[0m
				[33mIPs[0m:
				  [37mIP[0m:           [36m172.18.0.7[0m
				[33mControlled By[0m:  [36mReplicaSet/nginx[0m
				[33mContainers[0m:
				  [37mnginx[0m:
				    [33mContainer ID[0m:   [36mdocker://2885230a30908c8a6bda5a5366619c730b25b994eea61c931bba08ef4a8c8593[0m
				      [37mStarted[0m:      [36mSat, 10 Oct 2020 14:07:44 +0900[0m
			`),
		},
		{
			name:           "table format in kubectl describe can be colored by describe",
			darkBackground: true,
			tablePrinter:   NewTablePrinter(false, true, nil),
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
				[36m[0m  [32mType[0m             [35mStatus[0m  [37mLastHeartbeatTime[0m                 [33mLastTransitionTime[0m                [36mReason[0m                       [32mMessage[0m
				[36m[0m  [32m----[0m             [35m------[0m  [37m-----------------[0m                 [33m------------------[0m                [36m------[0m                       [32m-------[0m
				[36m[0m  [32mMemoryPressure[0m   [35mFalse[0m   [37mSun, 18 Oct 2020 12:00:54 +0900[0m   [33mWed, 14 Oct 2020 09:28:18 +0900[0m   [36mKubeletHasSufficientMemory[0m   [32mkubelet has sufficient memory available[0m
				[36m[0m  [32mDiskPressure[0m     [35mFalse[0m   [37mSun, 18 Oct 2020 12:00:54 +0900[0m   [33mWed, 14 Oct 2020 09:28:18 +0900[0m   [36mKubeletHasNoDiskPressure[0m     [32mkubelet has no disk pressure[0m
				[33mAddresses[0m:
				  [37mInternalIP[0m:  [36m172.17.0.3[0m
				  [37mHostname[0m:    [36mminikube[0m
				[33mCapacity[0m:
				  [37mcpu[0m:                [35m6[0m
				  [37mmemory[0m:             [36m2036900Ki[0m
				  [37mpods[0m:               [35m110[0m
				[33mAllocatable[0m:
				  [37mcpu[0m:                [35m6[0m
				  [37mmemory[0m:             [36m2036900Ki[0m
				  [37mpods[0m:               [35m110[0m
				[33mSystem Info[0m:
				  [37mMachine ID[0m:                 [36m55d2ccaefc9847c9a69356e7f3bd23f4[0m
				  [37mSystem UUID[0m:                [36mfe312784-2364-4bba-a55e-f56051539c21[0m
				[33mNon-terminated Pods[0m:          [36m(14 in total)[0m
				[36m[0m  [32mNamespace[0m                   [35mName[0m                                [37mCPU Requests[0m  [33mCPU Limits[0m  [36mMemory Requests[0m  [32mMemory Limits[0m  [32mAGE[0m
				[36m[0m  [32m---------[0m                   [35m----[0m                                [37m------------[0m  [33m----------[0m  [36m---------------[0m  [32m-------------[0m  [32m---[0m
				[36m[0m  [32mdefault[0m                     [35mnginx-6799fc88d8-dnmv5[0m              [37m0 (0%)[0m        [33m0 (0%)[0m      [36m0 (0%)[0m           [32m0 (0%)[0m         [32m7d21h[0m
				[36m[0m  [32mdefault[0m                     [35mnginx-6799fc88d8-m8pbc[0m              [37m0 (0%)[0m        [33m0 (0%)[0m      [36m0 (0%)[0m           [32m0 (0%)[0m         [32m7d21h[0m
				[36m[0m  [32mdefault[0m                     [35mnginx-6799fc88d8-qdf9b[0m              [37m0 (0%)[0m        [33m0 (0%)[0m      [36m0 (0%)[0m           [32m0 (0%)[0m         [32m7d21h[0m
				[33mAllocated resources[0m:
				  [36m(Total limits may be over 100 percent, i.e., overcommitted.)[0m
				[36m[0m  [32mResource[0m           [35mRequests[0m    [37mLimits[0m
				[36m[0m  [32m--------[0m           [35m--------[0m    [37m------[0m
				[36m[0m  [32mcpu[0m                [35m650m (10%)[0m  [37m0 (0%)[0m
				[36m[0m  [32mmemory[0m             [35m70Mi (3%)[0m   [37m170Mi (8%)[0m
				[33mEvents[0m:              [33m<none>[0m
			`),
		},
		{
			// This test input is invalid because contents in `Resource Quotas` have only 1 space as its indentation.
			// This is the bug of kubectl 1.19.3, and because of this
			// kubecolor was crashing.
			// A workaround PR is created for this and this test is making sure if it works.
			// So, this test should be removed after the kubectl fix
			// https://github.com/kubernetes/kubectl/issues/1005
			// For more details, see the PR description on GitHub.
			name:           "invalid test for the workaround",
			darkBackground: true,
			tablePrinter:   NewTablePrinter(false, true, nil),
			input: testutil.NewHereDoc(`
				Name:         default
				Labels:       <none>
				Annotations:  <none>
				Status:       Active
				
				Resource Quotas
				 Name:            mem-cpu-quota
				 Resource         Used  Hard
				 --------         ---   ---
				 limits.cpu       0     2
				 limits.memory    0     2Gi
				 requests.cpu     0     1
				 requests.memory  0     1Gi
				
				No LimitRange resource.`),
			expected: testutil.NewHereDoc(`
				[33mName[0m:         [36mdefault[0m
				[33mLabels[0m:       [33m<none>[0m
				[33mAnnotations[0m:  [33m<none>[0m
				[33mStatus[0m:       [36mActive[0m
				
				[36mResource Quotas[0m
				 [33mName[0m:            [36mmem-cpu-quota[0m
				[36m Resource[0m         [32mUsed[0m  [35mHard[0m
				[36m --------[0m         [32m---[0m   [35m---[0m
				[36m limits.cpu[0m       [32m0[0m     [35m2[0m
				[36m limits.memory[0m    [32m0[0m     [35m2Gi[0m
				[36m requests.cpu[0m     [32m0[0m     [35m1[0m
				[36m requests.memory[0m  [32m0[0m     [35m1Gi[0m
				
				[36mNo LimitRange resource.[0m
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
