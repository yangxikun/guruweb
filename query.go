package main

import (
	"bytes"
	"go/build"
	"go/token"
    "github.com/yangxikun/guruweb/internal/tools/guru"
    "github.com/yangxikun/guruweb/output"
    "strings"
)

func isMultiResMode(mode string) bool {
    if mode == "referrers" {
        return true
    }
    return false
}

// guru wrapper
func queryGuru(mode, pos string, jsonOut bool) (result []byte, err error) {
    guruQueryResult := make([][]byte, 0)
	jsonOutput := func(fset *token.FileSet, qr guru.QueryResult) {
		if jsonOut {
			// JSON output
            guruQueryResult = append(guruQueryResult, qr.JSON(fset))
		} else {
			// plain output
			out := &bytes.Buffer{}
			printf := func(pos interface{}, format string, args ...interface{}) {
				guru.Fprintf(out, fset, pos, format, args...)
			}
			qr.PrintPlain(printf)
            guruQueryResult = append(guruQueryResult, out.Bytes())
		}
	}

	ctxt := &build.Default
	query := guru.Query{
		Pos:        pos,
		Build:      ctxt,
		Scope:      defaultConfig.scope,
		PTALog:     nil,
		Reflection: false,
		Output:     jsonOutput,
	}

    if defaultConfig.verbose {
        output.Trace("guru -scope %s %s %s", strings.Join(defaultConfig.scope, ","), mode, pos)
    }

	err = guru.Run(mode, &query)
    if err != nil {
        return
    }

    if isMultiResMode(mode) {
        var _result bytes.Buffer
        if jsonOut {
            _result.WriteByte('[')
            for _, v := range guruQueryResult {
                _result.Write(v)
                _result.WriteByte(',')
            }
            _result.Truncate(_result.Len() - 1)
            _result.WriteByte(']')
        } else {
            for _, v := range guruQueryResult {
                _result.Write(v)
            }
        }
        result = _result.Bytes()
    } else {
        result = guruQueryResult[0]
    }
	return
}
