// Copyright (c) 2017 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
)

func TestFill(t *testing.T) {
	tests := [...]struct {
		name string
		src  string
		want string
	}{

		{
			name: "pointer to struct",
			src: `package p

import (
	"container/list"
	"io"
)

var b = myStruct{}
var s = myStruct{}

type myStruct struct {
	a *otherStruct
	b otherStruct
	c [1]*otherStruct
	d **otherStruct
}

type otherStruct struct {
	a list.Element
	b *list.Element
	io.Reader
}

type anotherStruct struct{ a int }`,
		},
	}

	for _, test := range tests {
		pkg, importNames, lit, typ := parseStruct(t, test.name, test.src)

		name := types.NewNamed(types.NewTypeName(0, pkg, "myStruct", nil), typ, nil)
		newlit, lines := zeroValue(pkg, importNames, lit, litInfo{typ: typ, name: name})

		out := printNode(t, test.name, newlit, lines)
		t.Log(out)
		//if test.want != out {
		//	t.Errorf("%q: got %v, want %v\n", test.name, out, test.want)
		//}
	}
}

func parseStruct(t *testing.T, filename, src string) (*types.Package, map[string]string, *ast.CompositeLit, *types.Struct) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	info := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	conf := types.Config{
		Importer: importer.Default(),
		Error:    func(err error) {},
	}

	pkg, _ := conf.Check(f.Name.Name, fset, []*ast.File{f}, &info)
	importNames := buildImportNameMap(f)

	expr := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Values[0]
	return pkg, importNames, expr.(*ast.CompositeLit), info.Types[expr].Type.Underlying().(*types.Struct)
}

func printNode(t *testing.T, name string, n ast.Node, lines int) string {
	fset := token.NewFileSet()
	file := fset.AddFile("", -1, lines)
	for i := 1; i <= lines; i++ {
		file.AddLine(i)
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, n); err != nil {
		t.Fatalf("%q: %v", name, err)
	}
	return buf.String()
}
