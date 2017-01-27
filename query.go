package guruweb

import (
	"bytes"
	"go/build"
	"go/token"
	"guruweb/internal/tools/guru"
)

func queryGuru(mode, pos string, json bool) (result []byte, err error) {
	output := func(fset *token.FileSet, qr guru.QueryResult) {
		result = qr.JSON(fset)
		if json {
			// JSON output
			result = qr.JSON(fset)
		} else {
			// plain output
			out := &bytes.Buffer{}
			printf := func(pos interface{}, format string, args ...interface{}) {
				guru.Fprintf(out, fset, pos, format, args...)
			}
			qr.PrintPlain(printf)
			result = out.Bytes()
		}
	}

	ctxt := &build.Default
	query := guru.Query{
		Pos:        pos,
		Build:      ctxt,
		Scope:      defaultConfig.scope,
		PTALog:     nil,
		Reflection: false,
		Output:     output,
	}

	err = guru.Run(mode, &query)
	return
}
