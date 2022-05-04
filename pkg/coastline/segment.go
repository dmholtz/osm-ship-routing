package coastline

type Segment interface {
	Size() int
	Left() int64
	Right() int64
	IsPolygon() bool
}

type AtomicSegment []int64

func NewAtomicSegment(raw []int64) AtomicSegment {
	if len(raw) < 2 {
		panic("len(AtomicSegment) < 2")
	}
	return AtomicSegment(raw)
}

func (as AtomicSegment) Size() int {
	return len(as)
}

func (as AtomicSegment) Left() int64 {
	return as[0]
}

func (as AtomicSegment) Right() int64 {
	return as[len(as)-1]
}

func (as AtomicSegment) IsPolygon() bool {
	return as.Left() == as.Right()
}

type ComposedSegment struct {
	size          int
	leftSegment   Segment // no pointers to interface: the interface either stores a struct or a reference to a struct
	leftInverted  bool
	left          int64
	rightSegment  Segment
	rightInverted bool
	right         int64
}

func NewComposedSegment(leftSeg Segment, leftInverted bool, rightSeg Segment, rightInverted bool) *ComposedSegment {
	size := leftSeg.Size() + rightSeg.Size() - 1

	left := leftSeg.Left()
	if leftInverted {
		left = leftSeg.Right()
	}
	right := rightSeg.Right()
	if rightInverted {
		right = rightSeg.Left()
	}
	return &ComposedSegment{size: size, leftSegment: leftSeg, leftInverted: leftInverted, left: left, rightSegment: rightSeg, rightInverted: rightInverted, right: right}
}

func (cs *ComposedSegment) Size() int {
	return cs.size
}

func (cs *ComposedSegment) Left() int64 {
	return cs.left
}

func (cs *ComposedSegment) Right() int64 {
	return cs.right
}

func (cs *ComposedSegment) IsPolygon() bool {
	return cs.Left() == cs.Right()
}

func (cs *ComposedSegment) LeftSegment() Segment {
	return cs.leftSegment
}

func (cs *ComposedSegment) RightSegment() Segment {
	return cs.rightSegment
}

func (cs *ComposedSegment) LeftInverted() bool {
	return cs.leftInverted
}

func (cs *ComposedSegment) RightInverted() bool {
	return cs.rightInverted
}
