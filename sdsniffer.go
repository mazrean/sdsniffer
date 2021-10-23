package sdsniffer

import (
	"context"
	"fmt"

	"github.com/mazrean/go-clone-detection"
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

	for _, clonePair := range clonePairs {
		pos := pass.Fset.Position(clonePair.Node2.Pos())
		end := pass.Fset.Position(clonePair.Node2.End())
		pass.ReportRangef(clonePair.Node1, "clone found: %s:%d:%d-%d:%d", pos.Filename, pos.Line, pos.Column, end.Line, end.Column)

		pos = pass.Fset.Position(clonePair.Node1.Pos())
		end = pass.Fset.Position(clonePair.Node1.End())
		pass.ReportRangef(clonePair.Node2, "clone found: %s:%d:%d-%d:%d", pos.Filename, pos.Line, pos.Column, end.Line, end.Column)
	}

	return nil, nil
}
