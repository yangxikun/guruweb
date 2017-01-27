// export guru identifiers, used by guruweb

package guru

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
)

// fprintf prints to w a message of the form "location: message\n"
// where location is derived from pos.
//
// pos must be one of:
//    - a token.Pos, denoting a position
//    - an ast.Node, denoting an interval
//    - anything with a Pos() method:
//         ssa.Member, ssa.Value, ssa.Instruction, types.Object, pointer.Label, etc.
//    - a QueryPos, denoting the extent of the user's query.
//    - nil, meaning no position at all.
//
// The output format is is compatible with the 'gnu'
// compilation-error-regexp in Emacs' compilation mode.
//
func Fprintf(w io.Writer, fset *token.FileSet, pos interface{}, format string, args ...interface{}) {
	var start, end token.Pos
	switch pos := pos.(type) {
	case ast.Node:
		start = pos.Pos()
		end = pos.End()
	case token.Pos:
		start = pos
		end = start
	case *types.PkgName:
		// The Pos of most PkgName objects does not coincide with an identifier,
		// so we suppress the usual start+len(name) heuristic for types.Objects.
		start = pos.Pos()
		end = start
	case types.Object:
		start = pos.Pos()
		end = start + token.Pos(len(pos.Name())) // heuristic
	case interface {
		Pos() token.Pos
	}:
		start = pos.Pos()
		end = start
	case *queryPos:
		start = pos.start
		end = pos.end
	case nil:
	// no-op
	default:
		panic(fmt.Sprintf("invalid pos: %T", pos))
	}

	if sp := fset.Position(start); start == end {
		// (prints "-: " for token.NoPos)
		fmt.Fprintf(w, "%s: ", sp)
	} else {
		ep := fset.Position(end)
		// The -1 below is a concession to Emacs's broken use of
		// inclusive (not half-open) intervals.
		// Other editors may not want it.
		// TODO(adonovan): add an -editor=vim|emacs|acme|auto
		// flag; auto uses EMACS=t / VIM=... / etc env vars.
		fmt.Fprintf(w, "%s:%d.%d-%d.%d: ",
			sp.Filename, sp.Line, sp.Column, ep.Line, ep.Column-1)
	}
	fmt.Fprintf(w, format, args...)
	io.WriteString(w, "\n")
}
