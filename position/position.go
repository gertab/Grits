package position

import "fmt"

type Position struct {
	StartLine int
	StartPos  int
	// EndLine   int
	// EndPos    int
}

func (p *Position) String() string {
	if p.StartLine == -1 && p.StartPos == -1 {
		// Position not set
		return ""
	}

	// return fmt.Sprintf("Line %d:%d", p.StartLine, p.StartPos)
	return fmt.Sprintf("Line %d", p.StartLine)
}

func (p *Position) New(start, end int) *Position {
	return &Position{start, end}
}

func (p *Position) Empty() *Position {
	return &Position{-1, -1}
}
