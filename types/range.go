package types

import "golang.org/x/tools/go/analysis"

type Range struct {
	analysis.Range
	fileName    string
	startLine   int
	startColumn int
	endLine     int
	endColumn   int
}

func NewRange(
	rng analysis.Range,
	fileName string,
	startLine int,
	startColumn int,
	endLine int,
	endColumn int,
) *Range {
	return &Range{
		Range:       rng,
		fileName:    fileName,
		startLine:   startLine,
		startColumn: startColumn,
		endLine:     endLine,
		endColumn:   endColumn,
	}
}

func (r *Range) FileName() string {
	return r.fileName
}

func (r *Range) StartLine() int {
	return r.startLine
}

func (r *Range) StartColumn() int {
	return r.startColumn
}

func (r *Range) EndLine() int {
	return r.endLine
}

func (r *Range) EndColumn() int {
	return r.endColumn
}
