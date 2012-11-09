package pretty

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/scanner"
	"go/token"
	"strings"
)

var replacer = strings.NewReplacer(
	",", ",\n",
	"{", "{\n",
	"{}", "{\n}",
	"}", "\n}",
)

const msg = "missing ',' before newline in composite literal"

func Print(x interface{}) string {
	dat := replacer.Replace(fmt.Sprintf("package foo;var x = %#v", x))
	set := token.NewFileSet()
	var (
		f   *ast.File
		err error
	)
	for {
		f, err = parser.ParseFile(set, "", dat, 0)
		if err == nil {
			break
		}

		serr, ok := err.(scanner.ErrorList)
		if !ok {
			panic(err)
		}
		changes := false
		bytes := []byte(dat)
		for i := len(serr) - 1; i >= 0; i-- {
			s := serr[i]
			if s.Msg != msg {
				continue
			}
			bytes = append(bytes[:s.Pos.Offset], append([]byte{','}, bytes[s.Pos.Offset:]...)...)
			changes = true
		}
		dat = string(bytes)
		if !changes {
			panic(err)
		}
	}

	var buf bytes.Buffer
	printer.Fprint(&buf, set, f)
	return buf.String()[21:]
}
