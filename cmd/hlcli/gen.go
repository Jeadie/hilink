// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	out := flag.String("o", "doc.go", "out file")
	pkg := flag.String("pkg", "github.com/jeadie/hilink", "go package")
	flag.Parse()
	if err := run(*out, *pkg); err != nil {
		log.Fatal(err)
	}
}

func run(out, pkg string) error {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, filepath.Join(os.Getenv("GOPATH"), "src", pkg), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	if len(pkgs) != 1 {
		return fmt.Errorf("invalid package count in %s", pkg)
	}
	// silly loop because it pkgs is a map ...
	var pkgName string
	for pkgName = range pkgs {
	}
	if pkgName != "hilink" {
		return fmt.Errorf("invalid package name %s", pkgName)
	}
	buf := new(bytes.Buffer)
	buf.WriteString(hdr)
	buf.WriteString("var methodParamMap = map[string][]string{\n")
	for _, f := range pkgs[pkgName].Files {
		for _, d := range f.Decls {
			fd, typ, ok := getRecvType(d)
			if !ok || typ != "Client" || !fd.Name.IsExported() || fd.Name.Name == "Do" {
				continue
			}
			str := `"` + fd.Name.Name + `": {`
			for _, p := range fd.Type.Params.List {
				for _, n := range p.Names {
					if n.Name == "ctx" {
						continue
					}
					str += `"` + n.Name + `",`
				}
			}
			str += "},\n"
			buf.WriteString(str)
		}
	}
	buf.WriteString("}\n\n")
	buf.WriteString("var methodCommentMap = map[string]string{\n")
	for _, f := range pkgs[pkgName].Files {
		for _, d := range f.Decls {
			fd, typ, ok := getRecvType(d)
			if !ok || typ != "Client" || !fd.Name.IsExported() || fd.Name.Name == "Do" {
				continue
			}
			str := `"` + fd.Name.Name + `": "` + strings.TrimSpace(strings.Replace(fd.Doc.Text(), "\n", " ", -1)) + "\",\n"
			buf.WriteString(str)
		}
	}
	buf.WriteString("}\n\n")
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(out, formatted, 0o644)
}

// getRecvType returns the receiver type.
func getRecvType(d ast.Decl) (*ast.FuncDecl, string, bool) {
	fd, ok := d.(*ast.FuncDecl)
	if !ok || fd.Recv == nil {
		return nil, "", false
	}
	se, ok := fd.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return nil, "", false
	}
	i, ok := se.X.(*ast.Ident)
	if !ok {
		return nil, "", false
	}
	return fd, i.Name, true
}

const (
	hdr = `package main

// Code generated by gen.go. DO NOT EDIT.

`
)
