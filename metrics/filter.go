package metrics

import "github.com/mazrean/sdsniffer/types"

type Filter struct {
	lineNumThreshold    int
	lineNumPerOperation int
}

func NewFilter(lineNumThreshold int, lineNumPerOperation int) *Filter {
	return &Filter{
		lineNumThreshold:    lineNumThreshold,
		lineNumPerOperation: lineNumPerOperation,
	}
}

func (f *Filter) Filter(rangePairs []*types.RangePair) ([]*types.RangePair, error) {
	// TODO: implement
	return rangePairs, nil
}
