package coastline

type Merger struct {
	workList   []Segment
	lookUpMap  map[int64]Segment
	polygons   []AtomicSegment
	mergeCount int
	unmergable int
}

func NewMerger(workList []Segment) *Merger {
	return &Merger{workList: workList, lookUpMap: make(map[int64]Segment), polygons: make([]AtomicSegment, 0)}
}

func (m *Merger) Merge() {

	for _, segment := range m.workList {

		left, lok, right, rok := m.lookUp(segment)
		for lok || rok {
			var match Segment
			if lok {
				match = left
			} else {
				match = right
			}

			segment = mergeSegmentsFast(segment, match)
			m.mergeCount++

			delete(m.lookUpMap, match.Left())
			delete(m.lookUpMap, match.Right())

			left, lok, right, rok = m.lookUp(segment)
		}

		if segment.IsPolygon() {
			polygon := createPolygon(segment)
			m.polygons = append(m.polygons, polygon)
		} else {
			m.lookUpMap[segment.Left()] = segment
			m.lookUpMap[segment.Right()] = segment
		}
	}
	m.unmergable = len(m.lookUpMap)
}

func (m Merger) Polygons() []AtomicSegment {
	return m.polygons
}

func (m Merger) MergeCount() int {
	return m.mergeCount
}

func (m Merger) UnmergableSegmentCount() int {
	return m.unmergable
}

func (m Merger) lookUp(s Segment) (Segment, bool, Segment, bool) {
	left, lok := m.lookUpMap[s.Left()]
	right, rok := m.lookUpMap[s.Right()]
	return left, lok, right, rok
}

func mergeSegmentsFast(first Segment, other Segment) *ComposedSegment {
	if first.Left() == other.Left() {
		// LeftInvertedRight
		return NewComposedSegment(first, true, other, false)
	} else if first.Left() == other.Right() {
		// RightLeft
		return NewComposedSegment(other, false, first, false)
	} else if first.Right() == other.Left() {
		// LeftRight
		return NewComposedSegment(first, false, other, false)
	} else if first.Right() == other.Right() {
		// LeftRightInverted
		return NewComposedSegment(first, false, other, true)
	} else {
		panic("Not mergeable!")
	}
}

func createPolygon(s Segment) AtomicSegment {
	// pre-allocate a slice with the final size of the polygon
	polygon := make(AtomicSegment, s.Size(), s.Size())
	createPolygonRecursive(s, false, &polygon, 0)
	return polygon
}

func createPolygonRecursive(s Segment, inverted bool, polygonPointer *AtomicSegment, idx int) {
	switch s := s.(type) {
	case AtomicSegment:
		if inverted {
			for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
				s[j], s[i] = s[i], s[j]
			}
		}
		copy((*polygonPointer)[idx:idx+s.Size()], s)
	case *ComposedSegment:
		if inverted {
			createPolygonRecursive(s.RightSegment(), !s.RightInverted(), polygonPointer, idx)
			// overwrite the overlapping item
			createPolygonRecursive(s.LeftSegment(), !s.LeftInverted(), polygonPointer, idx+s.RightSegment().Size()-1)
		} else {
			createPolygonRecursive(s.LeftSegment(), s.LeftInverted(), polygonPointer, idx)
			// overwrite the overlapping item
			createPolygonRecursive(s.RightSegment(), s.RightInverted(), polygonPointer, idx+s.LeftSegment().Size()-1)
		}
	default:
		panic("Segment is neither atomic nor composed.")
	}
}
