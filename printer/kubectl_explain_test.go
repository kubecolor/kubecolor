package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func Test_ExplainPrinter_Print(t *testing.T) {
	tests := []struct {
		name           string
		darkBackground bool
		recursive      bool
		input          string
		expected       string
	}{
		{
			name:           "kind, version, description, fields can be colorized with recursive=false",
			darkBackground: true,
			recursive:      false,
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
				
				   kind	<string>
				     Kind is a string value representing the REST resource this object
				     represents. Servers may infer this from the endpoint the client submits
				     requests to. Cannot be updated. In CamelCase. More info:
				     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
				
				   metadata	<Object>
				     Standard object's metadata. More info:
				     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
				
				   spec	<Object>
				     Spec defines the behavior of a node.
				     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
				
				   status	<Object>
				     Most recently observed status of the node. Populated by the system.
				     Read-only. More info:
				     https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`),
			expected: testutil.NewHereDoc(`
				[33mKIND[0m:     [36mNode[0m
				[33mVERSION[0m:  [36mv1[0m
				
				[33mDESCRIPTION[0m:
				     [36mNode is a worker node in Kubernetes. Each node will have a unique[0m
				     [36midentifier in the cache (i.e. in etcd).[0m
				
				[33mFIELDS[0m:
				   [37mapiVersion[0m	<[36mstring[0m>
				     [36mAPIVersion defines the versioned schema of this representation of an[0m
				     [36mobject. Servers should convert recognized schemas to the latest internal[0m
				     [36mvalue, and may reject unrecognized values. More info:[0m
				     [36mhttps://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources[0m
				
				   [37mkind[0m	<[36mstring[0m>
				     [36mKind is a string value representing the REST resource this object[0m
				     [36mrepresents. Servers may infer this from the endpoint the client submits[0m
				     [36mrequests to. Cannot be updated. In CamelCase. More info:[0m
				     [36mhttps://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds[0m
				
				   [37mmetadata[0m	<[36mObject[0m>
				     [36mStandard object's metadata. More info:[0m
				     [36mhttps://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata[0m
				
				   [37mspec[0m	<[36mObject[0m>
				     [36mSpec defines the behavior of a node.[0m
				     [36mhttps://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status[0m
				
				   [37mstatus[0m	<[36mObject[0m>
				     [36mMost recently observed status of the node. Populated by the system.[0m
				     [36mRead-only. More info:[0m
				     [36mhttps://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status[0m
			`),
		},
		{
			name:           "kind, version, description, fields can be colorized with recursive=true",
			darkBackground: true,
			recursive:      true,
			input: testutil.NewHereDoc(`
				KIND:     Node
				VERSION:  v1
				
				DESCRIPTION:
				     Node is a worker node in Kubernetes. Each node will have a unique
				     identifier in the cache (i.e. in etcd).
				
				FIELDS:
				   apiVersion	<string>
				   kind	<string>
				   metadata	<Object>
				      annotations	<map[string]string>
				      clusterName	<string>
				      creationTimestamp	<string>
				      deletionGracePeriodSeconds	<integer>
				      deletionTimestamp	<string>
				      finalizers	<[]string>
				      generateName	<string>
				      generation	<integer>
				      labels	<map[string]string>
				      managedFields	<[]Object>
				         apiVersion	<string>
				         fieldsType	<string>
				         fieldsV1	<map[string]>
				         manager	<string>
				         operation	<string>
				         time	<string>
				      name	<string>
				      namespace	<string>
				`),
			expected: testutil.NewHereDoc(`
				[33mKIND[0m:     [36mNode[0m
				[33mVERSION[0m:  [36mv1[0m
				
				[33mDESCRIPTION[0m:
				     [36mNode is a worker node in Kubernetes. Each node will have a unique[0m
				     [36midentifier in the cache (i.e. in etcd).[0m
				
				[33mFIELDS[0m:
				   [37mapiVersion[0m	<[36mstring[0m>
				   [37mkind[0m	<[36mstring[0m>
				   [37mmetadata[0m	<[36mObject[0m>
				      [37mannotations[0m	<[36mmap[string]string[0m>
				      [37mclusterName[0m	<[36mstring[0m>
				      [37mcreationTimestamp[0m	<[36mstring[0m>
				      [37mdeletionGracePeriodSeconds[0m	<[36minteger[0m>
				      [37mdeletionTimestamp[0m	<[36mstring[0m>
				      [37mfinalizers[0m	<[36m[]string[0m>
				      [37mgenerateName[0m	<[36mstring[0m>
				      [37mgeneration[0m	<[36minteger[0m>
				      [37mlabels[0m	<[36mmap[string]string[0m>
				      [37mmanagedFields[0m	<[36m[]Object[0m>
				         [33mapiVersion[0m	<[36mstring[0m>
				         [33mfieldsType[0m	<[36mstring[0m>
				         [33mfieldsV1[0m	<[36mmap[string][0m>
				         [33mmanager[0m	<[36mstring[0m>
				         [33moperation[0m	<[36mstring[0m>
				         [33mtime[0m	<[36mstring[0m>
				      [37mname[0m	<[36mstring[0m>
				      [37mnamespace[0m	<[36mstring[0m>
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := ExplainPrinter{DarkBackground: tt.darkBackground, Recursive: tt.recursive}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
