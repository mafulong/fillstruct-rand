package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var errNotFound = errors.New("no struct literal found at selection")
var W string

func main() {
	log.SetFlags(0)
	log.SetPrefix("fillstruct-rand: ")

	var (
		filename = flag.String("file", "", "required. filename")
		line     = flag.Int("line", 0, "required. line number of the struct literal")
		w        = flag.String("w", "", "optional. when set this, the generated code will write to the file which names w")
	)
	flag.Parse()

	if w != nil {
		W = *w
	}
	if (*line == 0) || *filename == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	path, err := absPath(*filename)
	if err != nil {
		log.Fatal(err)
	}

	var overlay map[string][]byte
	cfg := &packages.Config{
		Overlay: overlay,
		Mode:    packages.LoadAllSyntax,
		Tests:   true,
		Dir:     filepath.Dir(path),
		Fset:    token.NewFileSet(),
		Env:     os.Environ(),
	}

	pkgs, err := packages.Load(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if *line > 0 {
		err = byLine(pkgs, path, *line)
		switch err {
		case nil:
			return
		default:
			log.Fatal(err)
		}
	}
	log.Fatal(errNotFound)
}

func absPath(filename string) (string, error) {
	eval, err := filepath.EvalSymlinks(filename)
	if err != nil {
		return "", err
	}
	return filepath.Abs(eval)
}

func byLine(lprog []*packages.Package, path string, line int) (err error) {
	var f *ast.File
	var pkg *packages.Package
	for _, p := range lprog {
		for _, af := range p.Syntax {
			if file := p.Fset.File(af.Pos()); file.Name() == path {
				f = af
				pkg = p
			}
		}
	}
	if f == nil || pkg == nil {
		return fmt.Errorf("could not find file %q", path)
	}
	importNames := buildImportNameMap(f)
	rightStartLine := 0
	max := func(a, b int) int {
		if a > b {
			return a
		} else {
			return b
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		_ = lit
		startLine := pkg.Fset.Position(lit.Pos()).Line
		endLine := pkg.Fset.Position(lit.End()).Line
		if !(startLine <= line && line <= endLine) {
			return true
		}
		rightStartLine = max(startLine, rightStartLine)
		return true
	})

	var outs []output
	var prev types.Type
	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}
		startLine := pkg.Fset.Position(lit.Pos()).Line
		if !(startLine == rightStartLine) {
			return true
		}

		var info litInfo
		info.name, _ = pkg.TypesInfo.Types[lit].Type.(*types.Named)
		info.typ, ok = pkg.TypesInfo.Types[lit].Type.Underlying().(*types.Struct)
		if !ok {
			prev = pkg.TypesInfo.Types[lit].Type.Underlying()
			err = errNotFound
			return true
		}
		info.hideType = hideType(prev)

		startOff := pkg.Fset.Position(lit.Pos()).Offset
		endOff := pkg.Fset.Position(lit.End()).Offset
		newlit, lines := zeroValue(pkg.Types, importNames, lit, info)

		var out output
		out, err = prepareOutput(newlit, lines, startOff, endOff)
		if err != nil {
			return false
		}
		if Debug {
			log.Println(out.Code)
		} else if len(W) > 0 {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			contentStr := string(content)
			res := contentStr[:out.Start] + out.Code + contentStr[out.End:]
			ioutil.WriteFile(W, []byte(res), 0644)
			var resBytes []byte
			resBytes, err = imports.Process(path, []byte(res), nil)
			if err != nil {
				panic(err)
			}
			ioutil.WriteFile(W, resBytes, 0644)
		} else if len(W) == 0 {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			contentStr := string(content)
			res := contentStr[:out.Start] + out.Code + contentStr[out.End:]
			ioutil.WriteFile(path, []byte(res), 0644)
			var resBytes []byte
			resBytes, err = imports.Process(path, []byte(res), nil)
			if err != nil {
				log.Fatal("%+v", err)
				panic(err)
			}
			ioutil.WriteFile(path, resBytes, 0644)
		}

		outs = append(outs, out)
		return false
	})
	if err != nil {
		return err
	}
	if len(outs) == 0 {
		return errNotFound
	}

	for i := len(outs)/2 - 1; i >= 0; i-- {
		opp := len(outs) - 1 - i
		outs[i], outs[opp] = outs[opp], outs[i]
	}

	if err != nil {
		return err
	}
	//return json.NewEncoder(os.Stdout).Encode([]output{out})
	return nil
}

func hideType(t types.Type) bool {
	switch t.(type) {
	case *types.Array:
		return true
	case *types.Map:
		return true
	case *types.Slice:
		return true
	default:
		return false
	}
}

func buildImportNameMap(f *ast.File) map[string]string {
	imports := make(map[string]string)
	for _, i := range f.Imports {
		if i.Name != nil && i.Name.Name != "_" {
			path := i.Path.Value
			imports[path[1:len(path)-1]] = i.Name.Name
		}
	}
	return imports
}

type output struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Code  string `json:"code"`
}

func prepareOutput(n ast.Node, lines, start, end int) (output, error) {
	fset := token.NewFileSet()
	file := fset.AddFile("", -1, lines)
	for i := 1; i <= lines; i++ {
		file.AddLine(i)
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, n); err != nil {
		return output{}, err
	}
	return output{
		Start: start,
		End:   end,
		Code:  buf.String(),
	}, nil
}
