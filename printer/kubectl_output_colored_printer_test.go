package printer

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/kubecolor/kubecolor/kubectl"
	"github.com/kubecolor/kubecolor/testutil"
)

func Test_KubectlOutputColoredPrinter_Print(t *testing.T) {
	tests := []struct {
		name              string
		darkBackground    bool
		objFreshThreshold time.Duration
		subcommandInfo    *kubectl.SubcommandInfo
		input             string
		expected          string
	}{
		{
			name:           "kubectl top pod",
			darkBackground: true,
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
			name:           "kubectl top pod --no-headers",
			darkBackground: true,
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
			name:           "kubectl api-resources",
			darkBackground: true,
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
			name:           "kubectl api-resources --no-headers",
			darkBackground: true,
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
			name:           "kubectl get pod",
			darkBackground: true,
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
				[37mnginx-dnmv5[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-m8pbc[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-qdf9b[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
			`),
		},
		{
			name:           "kubectl get pod",
			darkBackground: true,
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
				[37mnginx-m8pbc[0m   [36m1/1[0m     [37mRunning[0m            [36m0[0m          [37m6d6h[0m
				[37mnginx-qdf9b[0m   [33m0/1[0m     [37mRunning[0m            [36m0[0m          [37m6d6h[0m
			`),
		},
		{
			name:              "kubectl get pod",
			darkBackground:    true,
			objFreshThreshold: time.Duration(300 * 1000000000),
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
				[37mnginx-dnmv6[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m6d6h[0m
				[37mnginx-m8pbc[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m5m[0m
				[37mnginx-qdf9b[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [32m4m59s[0m
			`),
		},
		{
			name:           "kubectl get pod --no-headers",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Get,
				NoHeader:   true,
			},
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
			name:           "kubectl get pod -o wide",
			darkBackground: true,
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
				[37mnginx-6799fc88d8-dnmv7[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m7d10h[0m   [36m172.18.0.5[0m   [37mminikube[0m   [36m<none>[0m           [37m<none>[0m
				[37mnginx-6799fc88d8-m8pbc[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m7d10h[0m   [36m172.18.0.4[0m   [37mminikube[0m   [36m<none>[0m           [37m<none>[0m
				[37mnginx-6799fc88d8-qdf9b[0m   [36m1/1[0m     [37mRunning[0m   [36m0[0m          [37m7d10h[0m   [36m172.18.0.3[0m   [37mminikube[0m   [36m<none>[0m           [37m<none>[0m
			`),
		},
		{
			name:           "kubectl get pod -o json",
			darkBackground: true,
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
			name:           "kubectl get pod -o yaml",
			darkBackground: true,
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
			name:           "kubectl describe pod",
			darkBackground: true,
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
				[33mLabels[0m:       [37mapp=nginx[0m
				[33mAnnotations[0m:  [33m<none>[0m
			`),
		},
		{
			name:           "kubectl api-versions",
			darkBackground: true,
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
			name:           "kubectl explain",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Explain,
			},
			input: testutil.NewHereDoc(`
				KIND:     Node
				VERSION:  v1

				DESCRIPTION:
				     Node is a worker node in Kubernetes. Each node will have a unique
				     identifier in the cache (i.e. in etcd).

				FIELDS:
				   apiVersion	<string>
				     APIVersion defines the versioned schema of this representation of an
				     object. Servers should convert recognized schemas to the latest internal
				     value, and may reject unrecognized values. More info:
				     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
				`),
			expected: testutil.NewHereDoc(`
				[33mKIND[0m:     [37mNode[0m
				[33mVERSION[0m:  [37mv1[0m

				[33mDESCRIPTION[0m:
				     [37mNode is a worker node in Kubernetes. Each node will have a unique[0m
				     [37midentifier in the cache (i.e. in etcd).[0m

				[33mFIELDS[0m:
				   [37mapiVersion[0m	<[37mstring[0m>
				     [37mAPIVersion defines the versioned schema of this representation of an[0m
				     [37mobject. Servers should convert recognized schemas to the latest internal[0m
				     [37mvalue, and may reject unrecognized values. More info:[0m
				     [37mhttps://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources[0m
			`),
		},
		{
			name:           "kubectl version",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Version,
			},
			input: testutil.NewHereDoc(`
				Client Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.3", GitCommit:"1e11e4a2108024935ecfcb2912226cedeafd99df", GitTreeState:"clean", BuildDate:"2020-10-14T18:49:28Z", GoVersion:"go1.15.2", Compiler:"gc", Platform:"darwin/amd64"}
				Server Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.2", GitCommit:"f5743093fd1c663cb0cbc89748f730662345d44d", GitTreeState:"clean", BuildDate:"2020-09-16T13:32:58Z", GoVersion:"go1.15", Compiler:"gc", Platform:"linux/amd64"}`),
			expected: testutil.NewHereDoc(`
				[33mClient Version[0m: [37mversion.Info[0m{[33mMajor[0m:"[37m1[0m", [33mMinor[0m:"[37m19[0m", [33mGitVersion[0m:"[37mv1.19.3[0m", [33mGitCommit[0m:"[37m1e11e4a2108024935ecfcb2912226cedeafd99df[0m", [33mGitTreeState[0m:"[37mclean[0m", [33mBuildDate[0m:"[37m2020-10-14T18:49:28Z[0m", [33mGoVersion[0m:"[37mgo1.15.2[0m", [33mCompiler[0m:"[37mgc[0m", [33mPlatform[0m:"[37mdarwin/amd64[0m"}
				[33mServer Version[0m: [37mversion.Info[0m{[33mMajor[0m:"[37m1[0m", [33mMinor[0m:"[37m19[0m", [33mGitVersion[0m:"[37mv1.19.2[0m", [33mGitCommit[0m:"[37mf5743093fd1c663cb0cbc89748f730662345d44d[0m", [33mGitTreeState[0m:"[37mclean[0m", [33mBuildDate[0m:"[37m2020-09-16T13:32:58Z[0m", [33mGoVersion[0m:"[37mgo1.15[0m", [33mCompiler[0m:"[37mgc[0m", [33mPlatform[0m:"[37mlinux/amd64[0m"}
			`),
		},
		{
			name:           "kubectl version --client",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Version,
			},
			input: testutil.NewHereDoc(`
				Client Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.3", GitCommit:"1e11e4a2108024935ecfcb2912226cedeafd99df", GitTreeState:"clean", BuildDate:"2020-10-14T18:49:28Z", GoVersion:"go1.15.2", Compiler:"gc", Platform:"darwin/amd64"}`),
			expected: testutil.NewHereDoc(`
				[33mClient Version[0m: [37mversion.Info[0m{[33mMajor[0m:"[37m1[0m", [33mMinor[0m:"[37m19[0m", [33mGitVersion[0m:"[37mv1.19.3[0m", [33mGitCommit[0m:"[37m1e11e4a2108024935ecfcb2912226cedeafd99df[0m", [33mGitTreeState[0m:"[37mclean[0m", [33mBuildDate[0m:"[37m2020-10-14T18:49:28Z[0m", [33mGoVersion[0m:"[37mgo1.15.2[0m", [33mCompiler[0m:"[37mgc[0m", [33mPlatform[0m:"[37mdarwin/amd64[0m"}
			`),
		},
		{
			name:           "kubectl version --short",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Version,
				Short:      true,
			},
			input: testutil.NewHereDoc(`
				Client Version: v1.19.3
				Server Version: v1.19.2`),
			expected: testutil.NewHereDoc(`
				[33mClient Version[0m: [37mv1.19.3[0m
				[33mServer Version[0m: [37mv1.19.2[0m
			`),
		},
		{
			name:           "kubectl version --short --client",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Version,
				Short:      true,
			},
			input: testutil.NewHereDoc(`
				Client Version: v1.19.3`),
			expected: testutil.NewHereDoc(`
				[33mClient Version[0m: [37mv1.19.3[0m
			`),
		},
		{
			name:           "kubectl options",
			darkBackground: true,
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
			name:           "kubectl apply -o json",
			darkBackground: true,
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
			name:           "kubectl apply -o yaml",
			darkBackground: true,
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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := KubectlOutputColoredPrinter{
				SubcommandInfo:    tt.subcommandInfo,
				DarkBackground:    tt.darkBackground,
				ObjFreshThreshold: tt.objFreshThreshold,
			}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
