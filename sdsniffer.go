package sdsniffer

import (
	"context"
	"fmt"

	"github.com/mazrean/go-clone-detection"
	"github.com/mazrean/sdsniffer/types"
	"golang.org/x/tools/go/analysis"
)

const doc = "sdsniffer is ..."

var (
	bufSize        int
	tokenThreshold int
)

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "sdsniffer",
	Doc:  doc,
	Run:  run,
}

func init() {
	Analyzer.Flags.IntVar(&bufSize, "buffer-size", 100, "buffer size")
	Analyzer.Flags.IntVar(&tokenThreshold, "token-threshold", 10, "Threshold for number of consecutive tokens")
}

func run(pass *analysis.Pass) (interface{}, error) {
	ctx := context.Background()

	cloneDetector := clone.NewCloneDetector(&clone.Config{
		BufSize:   bufSize,
		Threshold: tokenThreshold,
	})

	for _, file := range pass.Files {
		err := cloneDetector.AddNode(ctx, file)
		if err != nil {
			return nil, fmt.Errorf("failed to add node: %v", err)
		}
	}

	clonePairs, err := cloneDetector.GetClones()
	if err != nil {
		return nil, fmt.Errorf("failed to get clones: %v", err)
	}

	cloneRangePairs := make([]*types.RangePair, 0, len(clonePairs))
	for _, clonePair := range clonePairs {
		pos := pass.Fset.Position(clonePair.Node2.Pos())
		end := pass.Fset.Position(clonePair.Node2.End())
		rng1 := types.NewRange(
			clonePair.Node1,
			pos.Filename,
			pos.Line,
			pos.Column,
			end.Line,
			end.Column,
		)

		pos = pass.Fset.Position(clonePair.Node1.Pos())
		end = pass.Fset.Position(clonePair.Node1.End())
		rng2 := types.NewRange(
			clonePair.Node2,
			pos.Filename,
			pos.Line,
			pos.Column,
			end.Line,
			end.Column,
		)

		cloneRangePairs = append(cloneRangePairs, types.NewRangePair(rng1, rng2))
	}

	for _, cloneRangePair := range cloneRangePairs {
		rng1, rng2 := cloneRangePair.GetRanges()

		pass.ReportRangef(
			rng1,
			"clone found: %s:%d:%d-%d:%d",
			rng2.FileName(),
			rng2.StartLine(),
			rng2.StartColumn(),
			rng2.EndLine(),
			rng2.EndColumn(),
		)
		pass.ReportRangef(
			rng2,
			"clone found: %s:%d:%d-%d:%d",
			rng1.FileName(),
			rng1.StartLine(),
			rng1.StartColumn(),
			rng1.EndLine(),
			rng1.EndColumn(),
		)
	}

	return nil, nil
}
