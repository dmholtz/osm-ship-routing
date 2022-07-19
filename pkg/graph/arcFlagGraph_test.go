package graph

import (
	"math"
	"testing"
)

func TestReadFlag(t *testing.T) {
	t.Parallel()

	// test lower end
	fhe := FlaggedHalfEdge{To: 1, Weight: 1, Flag: 31}
	var p PartitionId
	for p = 0; p < 64; p++ {
		flagged := fhe.IsFlagged(p)
		expect := p < 5
		if flagged != expect {
			t.Errorf("Incorrect flag for partition %d, got: %t, want: %t.", p, flagged, expect)
		}
	}

	// test upper end
	fhe = FlaggedHalfEdge{To: 1, Weight: 1, Flag: 31 ^ math.MaxUint64}
	for p = 0; p < 63; p++ {
		flagged := fhe.IsFlagged(p)
		expect := p >= 5
		if flagged != expect {
			t.Errorf("Incorrect flag for partition %d, got: %t, want: %t.", p, flagged, expect)
		}
	}
}

func TestAddFlag(t *testing.T) {
	t.Parallel()

	fhe := FlaggedHalfEdge{To: 1, Weight: 1, Flag: 0}

	// test lower end
	var p PartitionId = 0
	fhe.AddFlag(p)
	if flagged := fhe.IsFlagged(p); !flagged {
		t.Errorf("Incorrect flag for partition %d, got: %t, want: %t.", p, flagged, true)
	}

	// overwrite existing flag
	fhe.AddFlag(p)
	if flagged := fhe.IsFlagged(p); !flagged {
		t.Errorf("Incorrect flag for partition %d, got: %t, want: %t.", p, flagged, true)
	}

	// test upper end
	p = 63
	fhe.AddFlag(p)
	if flagged := fhe.IsFlagged(p); !flagged {
		t.Errorf("Incorrect flag for partition %d, got: %t, want: %t.", p, flagged, true)
	}
}
