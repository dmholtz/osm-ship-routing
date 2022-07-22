package graph

import (
	"testing"
)

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

func TestComputeArcFlag2(t *testing.T) {
	t.Parallel()
	fg := NewFlaggedAdjacencyArrayFromFmi("out.fmi")

	fg1 := ComputeArcFlags(fg)

	t.Log(fg1.GetHalfEdgesFrom(0))

	WritePartitionedFmi(fg1, "out2.fmi")
}

func TestComputeArcFlag3(t *testing.T) {
	t.Parallel()
	fg := NewFlaggedAdjacencyArrayFromFmi("out0k.fmi")

	fg1 := ComputeArcFlags(fg)

	//t.Log(fg1.GetHalfEdgesFrom(0))

	WritePartitionedFmi(fg1, "out0k2.fmi")
}

func TestComputeArcFlag4(t *testing.T) {
	t.Parallel()
	fg := NewFlaggedAdjacencyArrayFromFmi(arcFlagGraphFile)

	fg1 := ComputeArcFlags(fg)

	WritePartitionedFmi(fg1, "out0k3.fmi")
}

func TestComputeArcFlag5(t *testing.T) {
	t.Parallel()
	fg := NewFlaggedAdjacencyArrayFromFmi(arcFlagGraphFile)

	fg1 := ComputeArcFlags(fg)

	WritePartitionedFmi(fg1, "out0k3.fmi")
}
