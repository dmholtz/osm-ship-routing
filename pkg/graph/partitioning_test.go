package graph

import (
	"testing"
)

func TestPartitioning(t *testing.T) {
	g := NewAdjacencyArrayFromFmi("../../graphs/ocean_10k.fmi")
	fg := GridPartitioning(g)
	WritePartitionedFmi(fg, "out.fmi")
}

func TestPartitioningFull(t *testing.T) {
	g := NewAdjacencyArrayFromFmi("../../graphs/ocean_equi_4.fmi")
	fg := GridPartitioning(g)
	WritePartitionedFmi(fg, "out0k.fmi")
}

func TestPartitioning0k(t *testing.T) {
	g := NewAdjacencyArrayFromFmi("../../graphs/ocean_0k.fmi")
	fg := GridPartitioning(g)
	WritePartitionedFmi(fg, "out0k.fmi")
}

func TestKdPartitioning(t *testing.T) {
	g := NewAdjacencyArrayFromFmi("testdata/kD_partitioning_3.fmi")
	fg := KdPartitioning(g, 3)
	WritePartitionedFmi(fg, "testdata/out.fmi")
}

func TestKdPartitioningFull(t *testing.T) {
	//g := NewAdjacencyArrayFromFmi("../../graphs/ocean_0k.fmi")
	g := NewAdjacencyArrayFromFmi("../../graphs/ocean_equi_4.fmi")
	fg := KdPartitioning(g, 6)
	WritePartitionedFmi(fg, "testdata/full.fmi")
}
