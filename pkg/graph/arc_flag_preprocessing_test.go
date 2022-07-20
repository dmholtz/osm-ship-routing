package graph

import "testing"

func TestComputeArcFlag(t *testing.T) {
	t.Parallel()
	fg := NewFlaggedAdjacencyArrayFromFmi(arcFlagGraphFile)

	fg1 := ComputeArcFlags(fg)

	for i := 0; i < fg1.NodeCount(); i++ {
		for _, e := range fg1.GetHalfEdgesFrom(i) {
			t.Logf("%d->%d: flag=%b", i, e.To, e.Flag)
		}
	}
}
