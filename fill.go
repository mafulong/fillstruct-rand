// Copyright (c) 2018 David R. Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"math/rand"
	"strconv"
	"strings"
)

var (
	// debug. special
	debug = false //debug使用
)

// litInfo contains the information about
// a literal to fill with zero values.
type litInfo struct {
	typ       types.Type   // the base type of the literal
	name      *types.Named // name of the type or nil, e.g. for an anonymous struct type
	hideType  bool         // flag to hide the element type inside an array, slice or map literal
	isPointer bool         // true if the literal is of a pointer type
}

type filler struct {
	pkg         *types.Package
	pos         token.Pos
	lines       int
	existing    map[string]*ast.KeyValueExpr
	first       bool
	importNames map[string]string // import path -> import name
}

// RandChoice random choice
func RandChoice(choices ...string) string {
	if len(choices) == 0 {
		return ""
	}
	i := rand.Intn(len(choices))
	return choices[i]
}

// RandomString  (n int, allowedChars ...[]rune)
func RandomString(n int, allowedChars ...[]rune) string {
	if n == 0 {
		n = rand.Intn(15)
	}
	var letters []rune

	if len(allowedChars) == 0 {
		letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// GetRandValueByType func
func GetRandValueByType(typeValue string, isStars ...bool) string {
	isStar := false
	if len(isStars) > 0 && isStars[0] {
		isStar = true
	}
	switch typeValue {
	case "context.Context":
		return "context.Background()"
	case "bool":
		if isStar {
			return fmt.Sprintf(`thrift.BoolPtr(%+v)`, GetRandValueByType(typeValue))
		}
		return RandChoice("true", "false")
	case "int", "int8", "int16", "uint",
		"uint8", "uint16":
		if isStar {
			return fmt.Sprintf(`thrift.IntPtr(int(%+v))`, GetRandValueByType(typeValue))
		}
		return ToStr(rand.Int31())
	case "int32":
		if isStar {
			return fmt.Sprintf(`thrift.Int32Ptr(int32(%+v))`, GetRandValueByType(typeValue))
		}
		return ToStr(rand.Int31())
	case "int64", "uint32", "uint64":
		if isStar {
			return fmt.Sprintf(`thrift.Int64Ptr(int64(%+v))`, GetRandValueByType(typeValue))
		}
		return ToStr(rand.Int63())
	case "string":
		if isStar {
			return fmt.Sprintf(`thrift.StringPtr(%+v)`, GetRandValueByType(typeValue))
		}
		return fmt.Sprintf("\"%+v\"", RandomString(10))
	case "float32":
		if isStar {
			return fmt.Sprintf(`thrift.Float32Ptr(%+v)`, GetRandValueByType(typeValue))
		}
		return fmt.Sprintf("%+v", rand.Float32())
	case "float64":
		if isStar {
			return fmt.Sprintf(`thrift.Float64Ptr(%+v)`, GetRandValueByType(typeValue))
		}
		return fmt.Sprintf("%+v", rand.Float64())
	}
	return ""
}

//GetRandValueByTypeV2  GetRandValueByTypeV2 return string
func GetRandValueByTypeV2(typeValue types.BasicKind, isStars ...bool) string {
	isStar := false
	if len(isStars) > 0 && isStars[0] {
		isStar = true
	}
	switch typeValue {
	case types.Bool:
		if isStar {
			return fmt.Sprintf(`thrift.BoolPtr(%+v)`, GetRandValueByTypeV2(typeValue))
		}
		return RandChoice("true", "false")
	case types.Int, types.Int8, types.Int16, types.Uint, types.Uint8, types.Uint16:
		if isStar {
			return fmt.Sprintf(`thrift.IntPtr(int(%+v))`, GetRandValueByTypeV2(typeValue))
		}
		return ToStr(rand.Int31())
	case types.Int32:
		if isStar {
			return fmt.Sprintf(`thrift.Int32Ptr(int32(%+v))`, GetRandValueByTypeV2(typeValue))
		}
		return ToStr(rand.Int31())
	case types.Int64, types.Uint64, types.Uint32:
		if isStar {
			return fmt.Sprintf(`thrift.Int64Ptr(int64(%+v))`, GetRandValueByTypeV2(typeValue))
		}
		return ToStr(rand.Int63())
	case types.String:
		if isStar {
			return fmt.Sprintf(`thrift.StringPtr(%+v)`, GetRandValueByTypeV2(typeValue))
		}
		return fmt.Sprintf("\"%+v\"", RandomString(10))
	case types.Float32:
		if isStar {
			return fmt.Sprintf(`thrift.Float32Ptr(%+v)`, GetRandValueByTypeV2(typeValue))
		}
		return fmt.Sprintf("%+v", rand.Float32())
	case types.Float64:
		if isStar {
			return fmt.Sprintf(`thrift.Float64Ptr(%+v)`, GetRandValueByTypeV2(typeValue))
		}
		return fmt.Sprintf("%+v", rand.Float64())
	}
	return ""
}

func zeroValue(pkg *types.Package, importNames map[string]string, lit *ast.CompositeLit, info litInfo) (ast.Expr, int) {
	f := filler{
		pkg:         pkg,
		pos:         1,
		first:       true,
		existing:    make(map[string]*ast.KeyValueExpr),
		importNames: importNames,
	}
	for _, e := range lit.Elts {
		kv := e.(*ast.KeyValueExpr)
		f.existing[kv.Key.(*ast.Ident).Name] = kv
	}
	return f.zero(info, make([]types.Type, 0, 8)), f.lines
}

func (f *filler) zero(info litInfo, visited []types.Type) ast.Expr {
	switch t := info.typ.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool, types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Float32, types.Float64, types.String:
			return &ast.Ident{Name: GetRandValueByTypeV2(t.Kind()), NamePos: f.pos}
		case types.Uintptr:
			return &ast.BasicLit{Value: "uintptr(0)", ValuePos: f.pos}
		case types.UnsafePointer:
			return &ast.BasicLit{Value: "unsafe.Pointer(uintptr(0))", ValuePos: f.pos}
		case types.Complex64, types.Complex128:
			return &ast.BasicLit{Value: "(0 + 0i)", ValuePos: f.pos}

		default:
			// Cannot create an expression for an invalid type.
			return nil
		}
	case *types.Chan:
		valTypeName, ok := typeString(f.pkg, f.importNames, t.Elem())
		if !ok {
			return nil
		}

		var dir ast.ChanDir
		switch t.Dir() {
		case types.SendRecv:
			dir = ast.SEND | ast.RECV
		case types.SendOnly:
			dir = ast.SEND
		case types.RecvOnly:
			dir = ast.RECV
		}

		return &ast.CallExpr{
			Fun: &ast.Ident{
				NamePos: f.pos,
				Name:    "make",
			},
			Lparen: f.pos,
			Args: []ast.Expr{
				&ast.ChanType{
					Dir:   dir,
					Value: ast.NewIdent(valTypeName),
				},
			},
			Rparen: f.pos,
		}
	case *types.Interface:
		if t.String() == "interface{Deadline() (deadline time.Time, ok bool); Done() <-chan struct{}; Err() error; Value(key interface{}) interface{}}" {
			return &ast.Ident{
				NamePos: f.pos,
				Name:    "context.Background()",
			}
		}
		return &ast.Ident{Name: "nil", NamePos: f.pos}
	case *types.Map:
		keyTypeName, ok := typeString(f.pkg, f.importNames, t.Key())
		if !ok {
			return nil
		}
		valTypeName, ok := typeString(f.pkg, f.importNames, t.Elem())
		if !ok {
			return nil
		}
		lit := &ast.CompositeLit{
			Lbrace: f.pos,
			Type: &ast.MapType{
				Map:   f.pos,
				Key:   ast.NewIdent(keyTypeName),
				Value: ast.NewIdent(valTypeName),
			},
		}
		f.pos++
		lit.Elts = []ast.Expr{
			&ast.KeyValueExpr{
				Key:   f.zero(litInfo{typ: t.Key(), name: info.name, hideType: true}, visited),
				Colon: f.pos,
				Value: f.zero(litInfo{typ: t.Elem(), name: info.name, hideType: true}, visited),
			},
		}
		f.pos++
		lit.Rbrace = f.pos
		f.lines += 2
		return lit
	case *types.Signature:
		params := make([]*ast.Field, t.Params().Len())
		for i := 0; i < t.Params().Len(); i++ {
			typeName, ok := typeString(f.pkg, f.importNames, t.Params().At(i).Type())
			if !ok {
				return nil
			}
			params[i] = &ast.Field{
				Type: ast.NewIdent(typeName),
			}
		}
		results := make([]*ast.Field, t.Results().Len())
		for i := 0; i < t.Results().Len(); i++ {
			typeName, ok := typeString(f.pkg, f.importNames, t.Results().At(i).Type())
			if !ok {
				return nil
			}
			results[i] = &ast.Field{
				Type: ast.NewIdent(typeName),
			}
		}
		return &ast.FuncLit{
			Type: &ast.FuncType{
				Func:    f.pos,
				Params:  &ast.FieldList{List: params},
				Results: &ast.FieldList{List: results},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{X: ast.NewIdent(`panic("not implemented")`)},
				},
			},
		}
	case *types.Slice:
		return &ast.Ident{Name: "nil", NamePos: f.pos}

	case *types.Array:
		lit := &ast.CompositeLit{Lbrace: f.pos}
		if !info.hideType {
			typeName, ok := typeString(f.pkg, f.importNames, t.Elem())
			if !ok {
				return nil
			}
			lit.Type = &ast.ArrayType{
				Lbrack: f.pos,
				Len:    &ast.BasicLit{Value: strconv.FormatInt(t.Len(), 10)},
				Elt:    ast.NewIdent(typeName),
			}
		}
		lit.Elts = make([]ast.Expr, 0, t.Len())
		for i := int64(0); i < t.Len(); i++ {
			f.pos++
			elemInfo := litInfo{typ: t.Elem().Underlying(), hideType: true}
			elemInfo.name, _ = t.Elem().(*types.Named)
			if v := f.zero(elemInfo, visited); v != nil {
				lit.Elts = append(lit.Elts, v)
			}
		}
		f.lines += len(lit.Elts) + 2
		f.pos++
		lit.Rbrace = f.pos
		return lit

	case *types.Named:
		if _, ok := t.Underlying().(*types.Struct); ok {
			info.name = t
		}
		info.typ = t.Underlying()
		return f.zero(info, visited)

	case *types.Pointer:
		if _, ok := t.Elem().Underlying().(*types.Struct); ok {
			info.typ = t.Elem()
			info.isPointer = true
			return f.zero(info, visited)
		}
		if _, ok := t.Elem().Underlying().(*types.Basic); ok {
			//fmt.Println(t.Elem().String(), t.Elem().Underlying(), t.Elem())
			info.typ = t.Elem()
			info.isPointer = true
			switch x := info.typ.(type) {
			case *types.Basic:
				switch x.Kind() {
				case types.Bool, types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Float32, types.Float64, types.String:
					return &ast.Ident{Name: GetRandValueByTypeV2(x.Kind(), true), NamePos: f.pos}
				}
			}
		}
		return &ast.Ident{Name: "nil", NamePos: f.pos}

	case *types.Struct:
		newlit := &ast.CompositeLit{Lbrace: f.pos}
		if !info.hideType && info.name != nil {
			typeName, ok := typeString(f.pkg, f.importNames, info.name)
			if !ok {
				return nil
			}
			newlit.Type = ast.NewIdent(typeName)
			if info.isPointer {
				newlit.Type.(*ast.Ident).Name = "&" + newlit.Type.(*ast.Ident).Name
			}
		} else if !info.hideType && info.name == nil {
			typeName, ok := typeString(f.pkg, f.importNames, t)
			if !ok {
				return nil
			}
			newlit.Type = ast.NewIdent(typeName)
		}

		for _, typ := range visited {
			if t == typ {
				return newlit
			}
		}
		visited = append(visited, t)

		first := f.first
		f.first = false
		lines := 0
		imported := isImported(f.pkg, info.name)

		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			// don't fill the field if it a gRPC system field
			if strings.HasPrefix(field.Name(), "XXX_") {
				continue
			}
			if kv, ok := f.existing[field.Name()]; first && ok {
				f.pos++
				lines++
				f.fixExprPos(kv)
				newlit.Elts = append(newlit.Elts, kv)
			} else if !ok && !imported || field.Exported() {
				f.pos++
				k := &ast.Ident{Name: field.Name(), NamePos: f.pos}
				if v := f.zero(litInfo{typ: field.Type(), name: nil}, visited); v != nil {
					lines++
					newlit.Elts = append(newlit.Elts, &ast.KeyValueExpr{
						Key:   k,
						Value: v,
					})
				} else {
					f.pos--
				}
			}
		}
		if lines > 0 {
			f.lines += lines + 2
			f.pos++
		}
		newlit.Rbrace = f.pos
		return newlit

	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (f *filler) fixExprPos(expr ast.Expr) {
	switch expr := expr.(type) {
	case nil:
		// ignore
	case *ast.BasicLit:
		expr.ValuePos = f.pos
	case *ast.BinaryExpr:
		f.fixExprPos(expr.X)
		expr.OpPos = f.pos
		f.fixExprPos(expr.Y)
	case *ast.CallExpr:
		f.fixExprPos(expr.Fun)
		expr.Lparen = f.pos
		for _, arg := range expr.Args {
			f.fixExprPos(arg)
		}
		expr.Rparen = f.pos
	case *ast.CompositeLit:
		f.fixExprPos(expr.Type)
		expr.Lbrace = f.pos
		for _, e := range expr.Elts {
			f.pos++
			f.fixExprPos(e)
		}
		if l := len(expr.Elts); l > 0 {
			f.lines += l + 2
		}
		f.pos++
		expr.Rbrace = f.pos
	case *ast.Ellipsis:
		expr.Ellipsis = f.pos
	case *ast.FuncLit:
		expr.Type.Func = f.pos
	case *ast.Ident:
		expr.NamePos = f.pos
	case *ast.IndexExpr:
		f.fixExprPos(expr.X)
		expr.Lbrack = f.pos
		f.fixExprPos(expr.Index)
		expr.Rbrack = f.pos
	case *ast.KeyValueExpr:
		f.fixExprPos(expr.Key)
		f.fixExprPos(expr.Value)
	case *ast.ParenExpr:
		expr.Lparen = f.pos
	case *ast.SelectorExpr:
		f.fixExprPos(expr.X)
		expr.Sel.NamePos = f.pos
	case *ast.SliceExpr:
		f.fixExprPos(expr.X)
		expr.Lbrack = f.pos
		f.fixExprPos(expr.Low)
		f.fixExprPos(expr.High)
		f.fixExprPos(expr.Max)
		expr.Rbrack = f.pos
	case *ast.StarExpr:
		expr.Star = f.pos
		f.fixExprPos(expr.X)
	case *ast.UnaryExpr:
		expr.OpPos = f.pos
		f.fixExprPos(expr.X)
	}
}

func isImported(pkg *types.Package, n *types.Named) bool {
	return n != nil && pkg != n.Obj().Pkg()
}
