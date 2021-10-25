package types

import (
	"errors"
	"go/token"
)

type Range struct {
	fileName    string
	lineNum     int
	pos         token.Pos
	startLine   int
	startColumn int
	end         token.Pos
	endLine     int
	endColumn   int
}

func NewRange(
	fileName string,
	pos token.Pos,
	startLine int,
	startColumn int,
	end token.Pos,
	endLine int,
	endColumn int,
) *Range {
	return &Range{
		fileName:    fileName,
		lineNum:     endLine - startLine + 1,
		pos:         pos,
		startLine:   startLine,
		startColumn: startColumn,
		end:         end,
		endLine:     endLine,
		endColumn:   endColumn,
	}
}

func (r *Range) FileName() string {
	return r.fileName
}

func (r *Range) Pos() token.Pos {
	return r.pos
}

func (r *Range) StartLine() int {
	return r.startLine
}

func (r *Range) StartColumn() int {
	return r.startColumn
}

func (r *Range) End() token.Pos {
	return r.end
}

func (r *Range) EndLine() int {
	return r.endLine
}

func (r *Range) EndColumn() int {
	return r.endColumn
}

func (r *Range) LineNum() int {
	return r.lineNum
}

func (r *Range) Join(r2 *Range) (*Range, error) {
	if r.FileName() != r2.FileName() {
		return nil, errors.New("cannot join ranges from different files")
	}

	var (
		pos         token.Pos
		startLine   int
		startColumn int
		end         token.Pos
		endLine     int
		endColumn   int
	)

	if r.StartLine() < r2.StartLine() {
		pos = r.pos
		startLine = r.StartLine()
		startColumn = r.StartColumn()
	} else if r.StartLine() > r2.StartLine() {
		pos = r2.pos
		startLine = r2.StartLine()
		startColumn = r2.StartColumn()
	} else {
		startLine = r.StartLine()
		if r.StartColumn() < r2.StartColumn() {
			pos = r.pos
			startColumn = r.StartColumn()
		} else {
			pos = r2.pos
			startColumn = r2.StartColumn()
		}
	}

	if r.EndLine() > r2.EndLine() {
		end = r.end
		endLine = r.EndLine()
		endColumn = r.EndColumn()
	} else if r.EndLine() < r2.EndLine() {
		end = r2.end
		endLine = r2.EndLine()
		endColumn = r2.EndColumn()
	} else {
		endLine = r.EndLine()
		if r.EndColumn() > r2.EndColumn() {
			end = r.end
			endColumn = r.EndColumn()
		} else {
			end = r2.end
			endColumn = r2.EndColumn()
		}
	}

	return &Range{
		fileName:    r.fileName,
		lineNum:     r.lineNum + r2.lineNum,
		pos:         pos,
		startLine:   startLine,
		startColumn: startColumn,
		end:         end,
		endLine:     endLine,
		endColumn:   endColumn,
	}, nil
}
