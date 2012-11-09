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

const (
	missingErr = "missing ',' before newline in composite literal"
	stringErr  = "string not terminated"
)

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
		byteDat := []byte(dat)
		for i := len(serr) - 1; i >= 0; i-- {
			s := serr[i]
			switch s.Msg {
			case missingErr:
				byteDat = append(byteDat[:s.Pos.Offset], append([]byte{','}, byteDat[s.Pos.Offset:]...)...)
				changes = true
			case stringErr:
				// these always come in pairs so do the one before it
				if i > 0 && serr[i-1].Msg == stringErr {
					i--
					s = serr[i]
				}
				index := bytes.IndexByte(byteDat[s.Pos.Offset:], '\n')
				if index == -1 {
					break
				}
				byteDat = append(byteDat[:s.Pos.Offset+index], byteDat[s.Pos.Offset+index+1:]...)
				changes = true
			}
		}
		dat = string(byteDat)
		if !changes {
			panic(err)
		}
	}

	var buf bytes.Buffer
	printer.Fprint(&buf, set, f)
	return strings.TrimRight(buf.String()[21:], "\n")
}
