package printer

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kubecolor/kubecolor/config"
	"github.com/kubecolor/kubecolor/config/testconfig"
	"github.com/kubecolor/kubecolor/kubectl"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_KubectlOutputColoredPrinter_Print(t *testing.T) {
	tests := []struct {
		name              string
		theme             *config.Theme
		objFreshThreshold time.Duration
		subcommandInfo    *kubectl.SubcommandInfo
		input             string
		expected          string
	}{
		{
			name:  "kubectl top pod",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Top,
			},
			input: testutil.NewHereDoc(`
				NAME        CPU(cores)   MEMORY(bytes)
				app-29twd   779m         221Mi
				app-2hhr6   1036m        220Mi
				app-52mbv   881m         137Mi`),
			expected: testutil.NewHereDoc(`
				[37mNAME        CPU(cores)   MEMORY(bytes)[0m
				[37mapp-29twd[0m   [36m779m[0m         [37m221Mi[0m
				[37mapp-2hhr6[0m   [36m1036m[0m        [37m220Mi[0m
				[37mapp-52mbv[0m   [36m881m[0m         [37m137Mi[0m
			`),
		},
		{
			name:  "kubectl top pod --no-headers",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Top,
				NoHeader:   true,
			},
			input: testutil.NewHereDoc(`
				app-29twd   779m         221Mi
				app-2hhr6   1036m        220Mi
				app-52mbv   881m         137Mi`),
			expected: testutil.NewHereDoc(`
				[37mapp-29twd[0m   [36m779m[0m         [37m221Mi[0m
				[37mapp-2hhr6[0m   [36m1036m[0m        [37m220Mi[0m
				[37mapp-52mbv[0m   [36m881m[0m         [37m137Mi[0m
			`),
		},
		{
			name:  "kubectl api-resources",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.APIResources,
			},
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
		{
			name:  "kubectl api-resources --no-headers",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.APIResources,
				NoHeader:   true,
			},
			input: testutil.NewHereDoc(`
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
		{
			name:  "kubectl get pod",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Get,
			},
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS    RESTARTS   AGE
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				[37mNAME          READY   STATUS    RESTARTS   AGE[0m
				[37mnginx-dnmv5[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-m8pbc[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-qdf9b[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m6d6h[0m
			`),
		},
		{
			name:  "kubectl get pod with crashloop",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Get,
			},
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS             RESTARTS   AGE
				nginx-dnmv4   1/1     CrashLoopBackOff   0          6d6h
				nginx-m8pbc   1/1     Running            0          6d6h
				nginx-qdf9b   0/1     Running            0          6d6h`),
			expected: testutil.NewHereDoc(`
				[37mNAME          READY   STATUS             RESTARTS   AGE[0m
				[37mnginx-dnmv4[0m   [36m1/1[0m     [31mCrashLoopBackOff[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-m8pbc[0m   [36m1/1[0m     [32mRunning[0m            [36m0[0m          [37m6d6h[0m
				[37mnginx-qdf9b[0m   [33m0/1[0m     [32mRunning[0m            [36m0[0m          [37m6d6h[0m
			`),
		},
		{
			name:              "kubectl get pod with fresh objects",
			theme:             testconfig.DarkTheme,
			objFreshThreshold: 5 * time.Minute,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Get,
			},
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS    RESTARTS   AGE
				nginx-dnmv6   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          5m
				nginx-qdf9b   1/1     Running   0          4m59s`),
			expected: testutil.NewHereDoc(`
				[37mNAME          READY   STATUS    RESTARTS   AGE[0m
				[37mnginx-dnmv6[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-m8pbc[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m5m[0m
				[37mnginx-qdf9b[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [32m4m59s[0m
			`),
		},
		{
			name:  "kubectl get pod --no-headers",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Get,
				NoHeader:   true,
			},
			input: testutil.NewHereDoc(`
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				[37mnginx-dnmv5[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-m8pbc[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-qdf9b[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m6d6h[0m
			`),
		},
		{
			name:  "kubectl get pod -o wide",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand:   kubectl.Get,
				FormatOption: kubectl.Wide,
			},
			input: testutil.NewHereDoc(`
				NAME                     READY   STATUS    RESTARTS   AGE     IP           NODE       NOMINATED NODE   READINESS GATES
				nginx-6799fc88d8-dnmv7   1/1     Running   0          7d10h   172.18.0.5   minikube   <none>           <none>
				nginx-6799fc88d8-m8pbc   1/1     Running   0          7d10h   172.18.0.4   minikube   <none>           <none>
				nginx-6799fc88d8-qdf9b   1/1     Running   0          7d10h   172.18.0.3   minikube   <none>           <none>`),
			expected: testutil.NewHereDoc(`
				[37mNAME                     READY   STATUS    RESTARTS   AGE     IP           NODE       NOMINATED NODE   READINESS GATES[0m
				[37mnginx-6799fc88d8-dnmv7[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m7d10h[0m   [36m172.18.0.5[0m   [37mminikube[0m   [36m<none>[0m           [37m<none>[0m
				[37mnginx-6799fc88d8-m8pbc[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m7d10h[0m   [36m172.18.0.4[0m   [37mminikube[0m   [36m<none>[0m           [37m<none>[0m
				[37mnginx-6799fc88d8-qdf9b[0m   [36m1/1[0m     [32mRunning[0m   [36m0[0m          [37m7d10h[0m   [36m172.18.0.3[0m   [37mminikube[0m   [36m<none>[0m           [37m<none>[0m
			`),
		},
		{
			name:  "kubectl get pod -o json",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand:   kubectl.Get,
				FormatOption: kubectl.Json,
			},
			input: testutil.NewHereDoc(`
				{
				    "apiVersion": "v1",
				    "kind": "Pod",
				    "num": 598,
				    "bool": true,
				    "null": null
				}`),
			expected: testutil.NewHereDoc(`
				{
				    "[37mapiVersion[0m": "[37mv1[0m",
				    "[37mkind[0m": "[37mPod[0m",
				    "[37mnum[0m": [35m598[0m,
				    "[37mbool[0m": [32mtrue[0m,
				    "[37mnull[0m": [33mnull[0m
				}
			`),
		},
		{
			name:  "kubectl get pod -o yaml",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand:   kubectl.Get,
				FormatOption: kubectl.Yaml,
			},
			input: testutil.NewHereDoc(`
				apiVersion: v1
				kind: "Pod"
				num: 415
				unknown: <unknown>
				none: <none>
				bool: true`),
			expected: testutil.NewHereDoc(`
				[33mapiVersion[0m: [37mv1[0m
				[33mkind[0m: "[37mPod[0m"
				[33mnum[0m: [35m415[0m
				[33munknown[0m: [33m<unknown>[0m
				[33mnone[0m: [33m<none>[0m
				[33mbool[0m: [32mtrue[0m
			`),
		},
		{
			name:  "kubectl describe pod",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Describe,
			},
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
				[33mName[0m:         [37mnginx-lpv5x[0m
				[33mNamespace[0m:    [37mdefault[0m
				[33mPriority[0m:     [35m0[0m
				[33mNode[0m:         [37mminikube/172.17.0.3[0m
				[33mReady[0m:        [32mtrue[0m
				[33mStart Time[0m:   [37mSat, 10 Oct 2020 14:07:17 +0900[0m
				[33mLabels[0m:       app=[37mnginx[0m
				[33mAnnotations[0m:  [33m<none>[0m
			`),
		},
		{
			name:  "kubectl api-versions",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.APIVersions,
			},
			input: testutil.NewHereDoc(`
				acme.cert-manager.io/v1alpha2
				admissionregistration.k8s.io/v1beta1
				apiextensions.k8s.io/v1beta1
				apiregistration.k8s.io/v1
				apiregistration.k8s.io/v1beta1
				apps/v1
				apps/v1beta1
				apps/v1beta2
				authentication.k8s.io/v1
				authentication.k8s.io/v1beta1
				authorization.k8s.io/v1
				authorization.k8s.io/v1beta1
				autoscaling/v1
				autoscaling/v2beta1
				autoscaling/v2beta2
				batch/v1
				batch/v1beta1`),
			expected: testutil.NewHereDoc(`
				[37macme.cert-manager.io/v1alpha2[0m
				[37madmissionregistration.k8s.io/v1beta1[0m
				[37mapiextensions.k8s.io/v1beta1[0m
				[37mapiregistration.k8s.io/v1[0m
				[37mapiregistration.k8s.io/v1beta1[0m
				[37mapps/v1[0m
				[37mapps/v1beta1[0m
				[37mapps/v1beta2[0m
				[37mauthentication.k8s.io/v1[0m
				[37mauthentication.k8s.io/v1beta1[0m
				[37mauthorization.k8s.io/v1[0m
				[37mauthorization.k8s.io/v1beta1[0m
				[37mautoscaling/v1[0m
				[37mautoscaling/v2beta1[0m
				[37mautoscaling/v2beta2[0m
				[37mbatch/v1[0m
				[37mbatch/v1beta1[0m
			`),
		},
		{
			name:  "kubectl version --client",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Version,
				Client:     true,
			},
			input: testutil.NewHereDoc(`
				Client Version: v1.29.0
				Kustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3
				`),
			expected: testutil.NewHereDoc(`
				[33mClient Version[0m: [37mv1.29.0[0m
				[33mKustomize Version[0m: [37mv5.0.4-0.20230601165947-6ce0bf390ce3[0m
				`),
		},
		{
			name:  "kubectl options",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Options,
			},
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
				[37mThe following options can be passed to any command:[0m

				      [33m--add-dir-header=false[0m: [37mIf true, adds the file directory to the header of the log messages[0m
				      [33m--alsologtostderr=false[0m: [37mlog to standard error as well as files[0m
				      [33m--as=''[0m: [37mUsername to impersonate for the operation[0m
				      [33m--as-group=[][0m: [37mGroup to impersonate for the operation, this flag can be repeated to specify multiple groups.[0m
				      [33m--cache-dir='/home/dtyler/.kube/cache'[0m: [37mDefault cache directory[0m
				      [33m--certificate-authority=''[0m: [37mPath to a cert file for the certificate authority[0m
				      [33m--client-certificate=''[0m: [37mPath to a client certificate file for TLS[0m
				      [33m--client-key=''[0m: [37mPath to a client key file for TLS[0m
				      [33m--cluster=''[0m: [37mThe name of the kubeconfig cluster to use[0m
				      [33m--context=''[0m: [37mThe name of the kubeconfig context to use[0m
				      [33m--insecure-skip-tls-verify=false[0m: [37mIf true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure[0m
			`),
		},
		{
			name:  "kubectl apply -o json",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand:   kubectl.Apply,
				FormatOption: kubectl.Json,
			},
			input: testutil.NewHereDoc(`
				{
				    "apiVersion": "apps/v1",
				    "kind": "Deployment",
				    "metadata": {
				        "annotations": {
				            "deployment.kubernetes.io/revision": "1",
				            "test": "false"
				        },
				        "creationTimestamp": "2020-11-04T13:14:07Z",
				        "generation": 3
				    }
				}`),
			expected: testutil.NewHereDoc(`
				{
				    "[37mapiVersion[0m": "[37mapps/v1[0m",
				    "[37mkind[0m": "[37mDeployment[0m",
				    "[37mmetadata[0m": {
				        "[33mannotations[0m": {
				            "[37mdeployment.kubernetes.io/revision[0m": "[37m1[0m",
				            "[37mtest[0m": "[37mfalse[0m"
				        },
				        "[33mcreationTimestamp[0m": "[37m2020-11-04T13:14:07Z[0m",
				        "[33mgeneration[0m": [35m3[0m
				    }
				}
			`),
		},
		{
			name:  "kubectl apply -o yaml",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand:   kubectl.Apply,
				FormatOption: kubectl.Yaml,
			},
			input: testutil.NewHereDoc(`
				apiVersion: apps/v1
				kind: Deployment
				metadata:
				  annotations:
				    deployment.kubernetes.io/revision: "1"
				    test: "false"
				  creationTimestamp: "2020-11-04T13:14:07Z"
				  generation: 3
				status:
				  availableReplicas: 3
				  conditions:
				  - lastTransitionTime: "2020-11-04T13:14:07Z"
				    lastUpdateTime: "2020-11-04T13:14:27Z"
				    message: ReplicaSet "nginx-f89759699" has successfully progressed.
				    reason: NewReplicaSetAvailable
				    status: "True"
				    type: Progressing
				  - lastTransitionTime: "2020-12-27T04:41:49Z"
				    lastUpdateTime: "2020-12-27T04:41:49Z"
				    message: Deployment has minimum availability.
				    reason: MinimumReplicasAvailable
				    status: "True"
				    type: Available
				  observedGeneration: 3
				  readyReplicas: 3
				  replicas: 3
				  updatedReplicas: 3
				`),
			expected: testutil.NewHereDoc(`
				[33mapiVersion[0m: [37mapps/v1[0m
				[33mkind[0m: [37mDeployment[0m
				[33mmetadata[0m:
				  [37mannotations[0m:
				    [33mdeployment.kubernetes.io/revision[0m: "[37m1[0m"
				    [33mtest[0m: "[37mfalse[0m"
				  [37mcreationTimestamp[0m: "[37m2020-11-04T13:14:07Z[0m"
				  [37mgeneration[0m: [35m3[0m
				[33mstatus[0m:
				  [37mavailableReplicas[0m: [35m3[0m
				  [37mconditions[0m:
				  - [33mlastTransitionTime[0m: "[37m2020-11-04T13:14:07Z[0m"
				    [33mlastUpdateTime[0m: "[37m2020-11-04T13:14:27Z[0m"
				    [33mmessage[0m: [37mReplicaSet "nginx-f89759699" has successfully progressed.[0m
				    [33mreason[0m: [37mNewReplicaSetAvailable[0m
				    [33mstatus[0m: "[37mTrue[0m"
				    [33mtype[0m: [37mProgressing[0m
				  - [33mlastTransitionTime[0m: "[37m2020-12-27T04:41:49Z[0m"
				    [33mlastUpdateTime[0m: "[37m2020-12-27T04:41:49Z[0m"
				    [33mmessage[0m: [37mDeployment has minimum availability.[0m
				    [33mreason[0m: [37mMinimumReplicasAvailable[0m
				    [33mstatus[0m: "[37mTrue[0m"
				    [33mtype[0m: [37mAvailable[0m
				  [37mobservedGeneration[0m: [35m3[0m
				  [37mreadyReplicas[0m: [35m3[0m
				  [37mreplicas[0m: [35m3[0m
				  [37mupdatedReplicas[0m: [35m3[0m
			`),
		},
		{
			name:  "kubectl events",
			theme: testconfig.DarkTheme,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Events,
			},
			input: testutil.NewHereDoc(`
				LAST SEEN   TYPE     REASON              OBJECT                       MESSAGE
				13s         Normal   SuccessfulCreate    ReplicaSet/nginx-76d6c9b8c   Created pod: nginx-76d6c9b8c-fmshc
				13s         Normal   SuccessfulCreate    ReplicaSet/nginx-76d6c9b8c   Created pod: nginx-76d6c9b8c-bkmwp
				13s         Normal   ScalingReplicaSet   Deployment/nginx             Scaled up replica set nginx-76d6c9b8c to 2
				12s         Normal   Scheduled           Pod/nginx-76d6c9b8c-fmshc    Successfully assigned default/nginx-76d6c9b8c-fmshc to minikube
				12s         Normal   Scheduled           Pod/nginx-76d6c9b8c-bkmwp    Successfully assigned default/nginx-76d6c9b8c-bkmwp to minikube
				12s         Normal   Pulling             Pod/nginx-76d6c9b8c-bkmwp    Pulling image "nginx"
				12s         Normal   Pulling             Pod/nginx-76d6c9b8c-fmshc    Pulling image "nginx"
				11s         Normal   Created             Pod/nginx-76d6c9b8c-bkmwp    Created container nginx
				11s         Normal   Started             Pod/nginx-76d6c9b8c-bkmwp    Started container nginx
				11s         Normal   Pulled              Pod/nginx-76d6c9b8c-bkmwp    Successfully pulled image "nginx" in 1.421388084s
				10s         Normal   Pulled              Pod/nginx-76d6c9b8c-fmshc    Successfully pulled image "nginx" in 2.892136877s
				9s          Normal   Created             Pod/nginx-76d6c9b8c-fmshc    Created container nginx
				9s          Normal   Started             Pod/nginx-76d6c9b8c-fmshc    Started container nginx
			`),
			expected: testutil.NewHereDoc(`
				[37mLAST SEEN   TYPE     REASON              OBJECT                       MESSAGE[0m
				[37m13s[0m         [32mNormal[0m   [32mSuccessfulCreate[0m    [36mReplicaSet/nginx-76d6c9b8c[0m   [37mCreated pod: nginx-76d6c9b8c-fmshc[0m
				[37m13s[0m         [32mNormal[0m   [32mSuccessfulCreate[0m    [36mReplicaSet/nginx-76d6c9b8c[0m   [37mCreated pod: nginx-76d6c9b8c-bkmwp[0m
				[37m13s[0m         [32mNormal[0m   [33mScalingReplicaSet[0m   [36mDeployment/nginx[0m             [37mScaled up replica set nginx-76d6c9b8c to 2[0m
				[37m12s[0m         [32mNormal[0m   [32mScheduled[0m           [36mPod/nginx-76d6c9b8c-fmshc[0m    [37mSuccessfully assigned default/nginx-76d6c9b8c-fmshc to minikube[0m
				[37m12s[0m         [32mNormal[0m   [32mScheduled[0m           [36mPod/nginx-76d6c9b8c-bkmwp[0m    [37mSuccessfully assigned default/nginx-76d6c9b8c-bkmwp to minikube[0m
				[37m12s[0m         [32mNormal[0m   [33mPulling[0m             [36mPod/nginx-76d6c9b8c-bkmwp[0m    [37mPulling image "nginx"[0m
				[37m12s[0m         [32mNormal[0m   [33mPulling[0m             [36mPod/nginx-76d6c9b8c-fmshc[0m    [37mPulling image "nginx"[0m
				[37m11s[0m         [32mNormal[0m   [32mCreated[0m             [36mPod/nginx-76d6c9b8c-bkmwp[0m    [37mCreated container nginx[0m
				[37m11s[0m         [32mNormal[0m   [32mStarted[0m             [36mPod/nginx-76d6c9b8c-bkmwp[0m    [37mStarted container nginx[0m
				[37m11s[0m         [32mNormal[0m   [32mPulled[0m              [36mPod/nginx-76d6c9b8c-bkmwp[0m    [37mSuccessfully pulled image "nginx" in 1.421388084s[0m
				[37m10s[0m         [32mNormal[0m   [32mPulled[0m              [36mPod/nginx-76d6c9b8c-fmshc[0m    [37mSuccessfully pulled image "nginx" in 2.892136877s[0m
				[37m9s[0m          [32mNormal[0m   [32mCreated[0m             [36mPod/nginx-76d6c9b8c-fmshc[0m    [37mCreated container nginx[0m
				[37m9s[0m          [32mNormal[0m   [32mStarted[0m             [36mPod/nginx-76d6c9b8c-fmshc[0m    [37mStarted container nginx[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := KubectlOutputColoredPrinter{
				SubcommandInfo:    tt.subcommandInfo,
				ObjFreshThreshold: tt.objFreshThreshold,
				Theme:             tt.theme,
			}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
