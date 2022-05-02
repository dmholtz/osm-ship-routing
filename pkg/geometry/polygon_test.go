package geometry

import (
	"testing"
)

func TestDummy(t *testing.T) {
	t.Parallel()
	if true != (4 > 3) {
		t.Errorf("want 'true', got '%t'", (4 > 3))
	}
}
