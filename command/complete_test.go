package command

import (
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestFilterCompleteResults(t *testing.T) {
	tests := []struct {
		name       string
		args       []CompleteArg
		toComplete string
		want       []CompleteArg
	}{
		{
			name:       "empty",
			args:       []CompleteArg{},
			toComplete: "",
			want:       []CompleteArg{},
		},

		{
			name: "filter flags by prefix",
			args: []CompleteArg{
				{Name: "--kubeconfig"},
				{Name: "--plain-foo"},
				{Name: "--plain-bar"},
				{Name: "--template"},
				{Name: "--watch"},
			},
			toComplete: "--plain",
			want: []CompleteArg{
				//{Name: "--kubeconfig"},
				{Name: "--plain-foo"},
				{Name: "--plain-bar"},
				//{Name: "--template"},
				//{Name: "--watch"},
			},
		},

		{
			name: "filter args by prefix",
			args: []CompleteArg{
				{Name: "default"},
				{Name: "kube-node-lease"},
				{Name: "kube-public"},
				{Name: "kube-system"},
				{Name: "metrics"},
			},
			toComplete: "kube",
			want: []CompleteArg{
				//{Name: "default"},
				{Name: "kube-node-lease"},
				{Name: "kube-public"},
				{Name: "kube-system"},
				//{Name: "metrics"},
			},
		},

		{
			// https://github.com/kubecolor/kubecolor/issues/207
			name: "filter flag value",
			args: []CompleteArg{
				{Name: "default"},
				{Name: "kube-node-lease"},
				{Name: "kube-public"},
				{Name: "kube-system"},
				{Name: "metrics"},
			},
			toComplete: "-n=kube",
			want: []CompleteArg{
				//{Name: "default"},
				{Name: "kube-node-lease"},
				{Name: "kube-public"},
				{Name: "kube-system"},
				//{Name: "metrics"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := filterCompleteResults(tc.args, tc.toComplete)
			testutil.MustEqual(t, tc.want, got)
		})
	}
}
