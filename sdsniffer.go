package sdsniffer

import (
	"context"
	"fmt"

	"github.com/mazrean/go-clone-detection"
	"golang.org/x/tools/go/analysis"
)

const doc = "sdsniffer is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "sdsniffer",
	Doc:  doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	ctx := context.Background()
	cloneDetector := clone.NewCloneDetector(nil)

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
