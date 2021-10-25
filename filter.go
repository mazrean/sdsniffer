package sdsniffer

import "github.com/mazrean/sdsniffer/types"

type Filter interface {
	Filter([]*types.RangePair) ([]*types.RangePair, error)
}

type NoFilter struct{}

func NewNoFilter() *NoFilter {
	return &NoFilter{}
}

func (n *NoFilter) Filter(rangePairs []*types.RangePair) ([]*types.RangePair, error) {
	return rangePairs, nil
}
