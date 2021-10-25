package types

type RangePair struct {
	range1 *Range
	range2 *Range
}

func NewRangePair(range1 *Range, range2 *Range) *RangePair {
	return &RangePair{range1, range2}
}

func (rp *RangePair) GetRanges() (*Range, *Range) {
	return rp.range1, rp.range2
}
