package metrics

import (
	"fmt"

	"github.com/mazrean/sdsniffer/types"
)

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
	filePairMap, err := f.createFilePairMap(rangePairs)
	if err != nil {
		return nil, fmt.Errorf("failed to create file pair map: %w", err)
	}

	rangePairs = []*types.RangePair{}
	for fileName1, rangePairMap := range filePairMap {
		for fileName2, rangePairList := range rangePairMap {
			newRangePairs, err := f.getRangePairs(fileName1 == fileName2, rangePairList)
			if err != nil {
				return nil, fmt.Errorf("failed to get range pairs: %w", err)
			}

			rangePairs = append(rangePairs, newRangePairs...)
		}
	}

	return rangePairs, nil
}

func (f *Filter) createFilePairMap(rangePairs []*types.RangePair) (map[string]map[string][]*types.RangePair, error) {
	filePairMap := map[string]map[string][]*types.RangePair{}
	for _, rangePair := range rangePairs {
		range1, range2 := rangePair.GetRanges()

		// ファイル名の組に対してmapの入る場所を1箇所に限定するため順番を入れ替える
		if range1.FileName() > range2.FileName() {
			rangePair.Swap()
			range1, range2 = rangePair.GetRanges()
		}

		if _, ok := filePairMap[range1.FileName()]; !ok {
			filePairMap[range1.FileName()] = map[string][]*types.RangePair{}
		}

		if _, ok := filePairMap[range1.FileName()][range2.FileName()]; !ok {
			filePairMap[range1.FileName()][range2.FileName()] = []*types.RangePair{}
		}

		filePairMap[range1.FileName()][range2.FileName()] = append(filePairMap[range1.FileName()][range2.FileName()], rangePair)
	}

	return filePairMap, nil
}

func (f *Filter) getRangePairs(isSameFile bool, rangePairs []*types.RangePair) ([]*types.RangePair, error) {
	filteredRangePairs := make([]*types.RangePair, 0)
	for _, rangePair := range rangePairs {
		if rangePair.LineNum() > f.lineNumThreshold {
			filteredRangePairs = append(filteredRangePairs, rangePair)
		}
	}

	chainedRangePairs := make([]*types.RangePair, len(filteredRangePairs))
	for i, rangePairI := range filteredRangePairs {
		if rangePairI == nil {
			if chainedRangePairs[i] == nil {
				continue
			} else {
				rangePairI = chainedRangePairs[i]
			}
		}

		for j, rangePairJ := range filteredRangePairs[i+1:] {
			if rangePairJ == nil {
				if chainedRangePairs[j] == nil {
					continue
				} else {
					rangePairJ = chainedRangePairs[j]
				}
			}

			i1, i2 := rangePairI.GetRanges()
			j1, j2 := rangePairJ.GetRanges()

			dist1, err := i1.Distance(j1)
			if err != nil {
				return nil, fmt.Errorf("failed to get distance: %w", err)
			}
			dist2, err := i2.Distance(j2)
			if err != nil {
				return nil, fmt.Errorf("failed to get distance: %w", err)
			}
			if dist1+dist2/2 <= f.lineNumThreshold {
				range1, err := i1.Join(j1)
				if err != nil {
					return nil, fmt.Errorf("failed to join range: %w", err)
				}
				range2, err := i2.Join(j2)
				if err != nil {
					return nil, fmt.Errorf("failed to join range: %w", err)
				}
				rangePair := types.NewRangePair(range1, range2)

				if rangePair.LineNum() > f.lineNumThreshold {
					chainedRangePairs = append(chainedRangePairs, rangePair)

					filteredRangePairs[i] = nil
					filteredRangePairs[j] = nil
					continue
				}
			}

			if isSameFile {
				dist1, err := i1.Distance(j2)
				if err != nil {
					return nil, fmt.Errorf("failed to get distance: %w", err)
				}
				dist2, err := i2.Distance(j1)
				if err != nil {
					return nil, fmt.Errorf("failed to get distance: %w", err)
				}
				if dist1+dist2/2 <= f.lineNumThreshold {
					range1, err := i1.Join(j2)
					if err != nil {
						return nil, fmt.Errorf("failed to join range: %w", err)
					}
					range2, err := i2.Join(j1)
					if err != nil {
						return nil, fmt.Errorf("failed to join range: %w", err)
					}
					rangePair := types.NewRangePair(range1, range2)

					if rangePair.LineNum() > f.lineNumThreshold {
						chainedRangePairs = append(chainedRangePairs, rangePair)

						filteredRangePairs[i] = nil
						filteredRangePairs[j] = nil
					}
				}
			}
		}
	}

	newRangePairs := make([]*types.RangePair, 0)
	for _, chainedRangePairs := range chainedRangePairs {
		if chainedRangePairs != nil {
			newRangePairs = append(newRangePairs, chainedRangePairs)
		}
	}

	for _, rangePair := range filteredRangePairs {
		if rangePair != nil {
			if rangePair.LineNum() > f.lineNumPerOperation {
				newRangePairs = append(newRangePairs, rangePair)
			}
		}
	}

	return newRangePairs, nil
}
