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
				\e[37mNAME        CPU(cores)   MEMORY(bytes)\e[0m
				\e[37mapp-29twd\e[0m   \e[36m779m\e[0m         \e[37m221Mi\e[0m
				\e[37mapp-2hhr6\e[0m   \e[36m1036m\e[0m        \e[37m220Mi\e[0m
				\e[37mapp-52mbv\e[0m   \e[36m881m\e[0m         \e[37m137Mi\e[0m
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
				\e[37mapp-29twd\e[0m   \e[36m779m\e[0m         \e[37m221Mi\e[0m
				\e[37mapp-2hhr6\e[0m   \e[36m1036m\e[0m        \e[37m220Mi\e[0m
				\e[37mapp-52mbv\e[0m   \e[36m881m\e[0m         \e[37m137Mi\e[0m
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
				\e[37mNAME          READY   STATUS    RESTARTS   AGE\e[0m
				\e[37mnginx-dnmv5\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-m8pbc\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-qdf9b\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
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
				\e[37mNAME          READY   STATUS             RESTARTS   AGE\e[0m
				\e[37mnginx-dnmv4\e[0m   \e[36m1/1\e[0m     \e[31mCrashLoopBackOff\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-m8pbc\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m            \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-qdf9b\e[0m   \e[33m0/1\e[0m     \e[32mRunning\e[0m            \e[36m0\e[0m          \e[37m6d6h\e[0m
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
				\e[37mNAME          READY   STATUS    RESTARTS   AGE\e[0m
				\e[37mnginx-dnmv6\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-m8pbc\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m5m\e[0m
				\e[37mnginx-qdf9b\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[32m4m59s\e[0m
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
				\e[37mnginx-dnmv5\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-m8pbc\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
				\e[37mnginx-qdf9b\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m6d6h\e[0m
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
				\e[37mNAME                     READY   STATUS    RESTARTS   AGE     IP           NODE       NOMINATED NODE   READINESS GATES\e[0m
				\e[37mnginx-6799fc88d8-dnmv7\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m7d10h\e[0m   \e[36m172.18.0.5\e[0m   \e[37mminikube\e[0m   \e[36m<none>\e[0m           \e[37m<none>\e[0m
				\e[37mnginx-6799fc88d8-m8pbc\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m7d10h\e[0m   \e[36m172.18.0.4\e[0m   \e[37mminikube\e[0m   \e[36m<none>\e[0m           \e[37m<none>\e[0m
				\e[37mnginx-6799fc88d8-qdf9b\e[0m   \e[36m1/1\e[0m     \e[32mRunning\e[0m   \e[36m0\e[0m          \e[37m7d10h\e[0m   \e[36m172.18.0.3\e[0m   \e[37mminikube\e[0m   \e[36m<none>\e[0m           \e[37m<none>\e[0m
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
				    "\e[37mapiVersion\e[0m": "\e[37mv1\e[0m",
				    "\e[37mkind\e[0m": "\e[37mPod\e[0m",
				    "\e[37mnum\e[0m": \e[35m598\e[0m,
				    "\e[37mbool\e[0m": \e[32mtrue\e[0m,
				    "\e[37mnull\e[0m": \e[33mnull\e[0m
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
				\e[33mapiVersion\e[0m: \e[37mv1\e[0m
				\e[33mkind\e[0m: "\e[37mPod\e[0m"
				\e[33mnum\e[0m: \e[35m415\e[0m
				\e[33munknown\e[0m: \e[33m<unknown>\e[0m
				\e[33mnone\e[0m: \e[33m<none>\e[0m
				\e[33mbool\e[0m: \e[32mtrue\e[0m
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
				\e[33mName\e[0m:         \e[37mnginx-lpv5x\e[0m
				\e[33mNamespace\e[0m:    \e[37mdefault\e[0m
				\e[33mPriority\e[0m:     \e[35m0\e[0m
				\e[33mNode\e[0m:         \e[37mminikube/172.17.0.3\e[0m
				\e[33mReady\e[0m:        \e[32mtrue\e[0m
				\e[33mStart Time\e[0m:   \e[37mSat, 10 Oct 2020 14:07:17 +0900\e[0m
				\e[33mLabels\e[0m:       app=\e[37mnginx\e[0m
				\e[33mAnnotations\e[0m:  \e[33m<none>\e[0m
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
				\e[37macme.cert-manager.io/v1alpha2\e[0m
				\e[37madmissionregistration.k8s.io/v1beta1\e[0m
				\e[37mapiextensions.k8s.io/v1beta1\e[0m
				\e[37mapiregistration.k8s.io/v1\e[0m
				\e[37mapiregistration.k8s.io/v1beta1\e[0m
				\e[37mapps/v1\e[0m
				\e[37mapps/v1beta1\e[0m
				\e[37mapps/v1beta2\e[0m
				\e[37mauthentication.k8s.io/v1\e[0m
				\e[37mauthentication.k8s.io/v1beta1\e[0m
				\e[37mauthorization.k8s.io/v1\e[0m
				\e[37mauthorization.k8s.io/v1beta1\e[0m
				\e[37mautoscaling/v1\e[0m
				\e[37mautoscaling/v2beta1\e[0m
				\e[37mautoscaling/v2beta2\e[0m
				\e[37mbatch/v1\e[0m
				\e[37mbatch/v1beta1\e[0m
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
				\e[33mClient Version\e[0m: \e[37mv1.29.0\e[0m
				\e[33mKustomize Version\e[0m: \e[37mv5.0.4-0.20230601165947-6ce0bf390ce3\e[0m
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
				    "\e[37mapiVersion\e[0m": "\e[37mapps/v1\e[0m",
				    "\e[37mkind\e[0m": "\e[37mDeployment\e[0m",
				    "\e[37mmetadata\e[0m": {
				        "\e[33mannotations\e[0m": {
				            "\e[37mdeployment.kubernetes.io/revision\e[0m": "\e[37m1\e[0m",
				            "\e[37mtest\e[0m": "\e[37mfalse\e[0m"
				        },
				        "\e[33mcreationTimestamp\e[0m": "\e[37m2020-11-04T13:14:07Z\e[0m",
				        "\e[33mgeneration\e[0m": \e[35m3\e[0m
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
				\e[33mapiVersion\e[0m: \e[37mapps/v1\e[0m
				\e[33mkind\e[0m: \e[37mDeployment\e[0m
				\e[33mmetadata\e[0m:
				  \e[37mannotations\e[0m:
				    \e[33mdeployment.kubernetes.io/revision\e[0m: "\e[37m1\e[0m"
				    \e[33mtest\e[0m: "\e[37mfalse\e[0m"
				  \e[37mcreationTimestamp\e[0m: "\e[37m2020-11-04T13:14:07Z\e[0m"
				  \e[37mgeneration\e[0m: \e[35m3\e[0m
				\e[33mstatus\e[0m:
				  \e[37mavailableReplicas\e[0m: \e[35m3\e[0m
				  \e[37mconditions\e[0m:
				  - \e[33mlastTransitionTime\e[0m: "\e[37m2020-11-04T13:14:07Z\e[0m"
				    \e[33mlastUpdateTime\e[0m: "\e[37m2020-11-04T13:14:27Z\e[0m"
				    \e[33mmessage\e[0m: \e[37mReplicaSet "nginx-f89759699" has successfully progressed.\e[0m
				    \e[33mreason\e[0m: \e[37mNewReplicaSetAvailable\e[0m
				    \e[33mstatus\e[0m: "\e[37mTrue\e[0m"
				    \e[33mtype\e[0m: \e[37mProgressing\e[0m
				  - \e[33mlastTransitionTime\e[0m: "\e[37m2020-12-27T04:41:49Z\e[0m"
				    \e[33mlastUpdateTime\e[0m: "\e[37m2020-12-27T04:41:49Z\e[0m"
				    \e[33mmessage\e[0m: \e[37mDeployment has minimum availability.\e[0m
				    \e[33mreason\e[0m: \e[37mMinimumReplicasAvailable\e[0m
				    \e[33mstatus\e[0m: "\e[37mTrue\e[0m"
				    \e[33mtype\e[0m: \e[37mAvailable\e[0m
				  \e[37mobservedGeneration\e[0m: \e[35m3\e[0m
				  \e[37mreadyReplicas\e[0m: \e[35m3\e[0m
				  \e[37mreplicas\e[0m: \e[35m3\e[0m
				  \e[37mupdatedReplicas\e[0m: \e[35m3\e[0m
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
				\e[37mLAST SEEN   TYPE     REASON              OBJECT                       MESSAGE\e[0m
				\e[37m13s\e[0m         \e[32mNormal\e[0m   \e[32mSuccessfulCreate\e[0m    \e[36mReplicaSet/nginx-76d6c9b8c\e[0m   \e[37mCreated pod: nginx-76d6c9b8c-fmshc\e[0m
				\e[37m13s\e[0m         \e[32mNormal\e[0m   \e[32mSuccessfulCreate\e[0m    \e[36mReplicaSet/nginx-76d6c9b8c\e[0m   \e[37mCreated pod: nginx-76d6c9b8c-bkmwp\e[0m
				\e[37m13s\e[0m         \e[32mNormal\e[0m   \e[33mScalingReplicaSet\e[0m   \e[36mDeployment/nginx\e[0m             \e[37mScaled up replica set nginx-76d6c9b8c to 2\e[0m
				\e[37m12s\e[0m         \e[32mNormal\e[0m   \e[32mScheduled\e[0m           \e[36mPod/nginx-76d6c9b8c-fmshc\e[0m    \e[37mSuccessfully assigned default/nginx-76d6c9b8c-fmshc to minikube\e[0m
				\e[37m12s\e[0m         \e[32mNormal\e[0m   \e[32mScheduled\e[0m           \e[36mPod/nginx-76d6c9b8c-bkmwp\e[0m    \e[37mSuccessfully assigned default/nginx-76d6c9b8c-bkmwp to minikube\e[0m
				\e[37m12s\e[0m         \e[32mNormal\e[0m   \e[33mPulling\e[0m             \e[36mPod/nginx-76d6c9b8c-bkmwp\e[0m    \e[37mPulling image "nginx"\e[0m
				\e[37m12s\e[0m         \e[32mNormal\e[0m   \e[33mPulling\e[0m             \e[36mPod/nginx-76d6c9b8c-fmshc\e[0m    \e[37mPulling image "nginx"\e[0m
				\e[37m11s\e[0m         \e[32mNormal\e[0m   \e[32mCreated\e[0m             \e[36mPod/nginx-76d6c9b8c-bkmwp\e[0m    \e[37mCreated container nginx\e[0m
				\e[37m11s\e[0m         \e[32mNormal\e[0m   \e[32mStarted\e[0m             \e[36mPod/nginx-76d6c9b8c-bkmwp\e[0m    \e[37mStarted container nginx\e[0m
				\e[37m11s\e[0m         \e[32mNormal\e[0m   \e[32mPulled\e[0m              \e[36mPod/nginx-76d6c9b8c-bkmwp\e[0m    \e[37mSuccessfully pulled image "nginx" in 1.421388084s\e[0m
				\e[37m10s\e[0m         \e[32mNormal\e[0m   \e[32mPulled\e[0m              \e[36mPod/nginx-76d6c9b8c-fmshc\e[0m    \e[37mSuccessfully pulled image "nginx" in 2.892136877s\e[0m
				\e[37m9s\e[0m          \e[32mNormal\e[0m   \e[32mCreated\e[0m             \e[36mPod/nginx-76d6c9b8c-fmshc\e[0m    \e[37mCreated container nginx\e[0m
				\e[37m9s\e[0m          \e[32mNormal\e[0m   \e[32mStarted\e[0m             \e[36mPod/nginx-76d6c9b8c-fmshc\e[0m    \e[37mStarted container nginx\e[0m
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
