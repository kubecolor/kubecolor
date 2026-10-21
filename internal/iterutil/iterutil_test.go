package iterutil

import (
	"slices"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestWindowedPair(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		seq := slices.Values([]int{})
		got := slices.Collect(WindowedPair(seq))

		if len(got) != 0 {
			t.Fatalf("expected no results, but got: %#v", got)
		}
	})

	t.Run("single", func(t *testing.T) {
		seq := slices.Values([]int{1})
		got := slices.Collect(WindowedPair(seq))
		want := []NextPair[int]{
			{Curr: 1, Next: 0, HasNext: false},
		}

		testutil.Equal(t, want, got)
	})

	t.Run("multiple", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		got := slices.Collect(WindowedPair(seq))
		want := []NextPair[int]{
			{Curr: 1, Next: 2, HasNext: true},
			{Curr: 2, Next: 3, HasNext: true},
			{Curr: 3, Next: 0, HasNext: false},
		}

		testutil.Equal(t, want, got)
	})
}

func TestWindowedPair2(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		seq := slices.All([]int{})
		got := slices.Collect(WindowedPair2(seq))

		if len(got) != 0 {
			t.Fatalf("expected no results, but got: %#v", got)
		}
	})

	t.Run("single", func(t *testing.T) {
		seq := slices.All([]int{1})
		got := slices.Collect(WindowedPair2(seq))
		want := []NextPair2[int, int]{
			{Curr0: 0, Curr1: 1, Next0: 0, Next1: 0, HasNext: false},
		}

		testutil.Equal(t, want, got)
	})

	t.Run("multiple", func(t *testing.T) {
		seq := slices.All([]int{1, 2, 3})
		got := slices.Collect(WindowedPair2(seq))
		want := []NextPair2[int, int]{
			{Curr0: 0, Curr1: 1, Next0: 1, Next1: 2, HasNext: true},
			{Curr0: 1, Curr1: 2, Next0: 2, Next1: 3, HasNext: true},
			{Curr0: 2, Curr1: 3, Next0: 0, Next1: 0, HasNext: false},
		}

		testutil.Equal(t, want, got)
	})
}
