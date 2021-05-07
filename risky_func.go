package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/thoas/go-funk"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"golang.org/x/tools/cover"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func main() {
	profiles, err := cover.ParseProfiles(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	var total, covered int64

	var funcInfos []*funcInfo
	for _, profile := range profiles {
		fn := profile.FileName
		file, err := findFile(fn)
		if err != nil {
			fmt.Println(err)
			return
		}
		funcs, err := findFuncs(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Now match up functions and profile blocks.
		for _, f := range funcs {
			c, t := f.coverage(profile)
			funcInfos = append(funcInfos,
				&funcInfo{fileName: fn,
					functionName:      f.name,
					functionStartLine: f.startLine,
					uncoveredLines:    t - c})
			total += t
			covered += c
		}
	}

	sort.Slice(funcInfos, func(i, j int) bool {
		return funcInfos[i].uncoveredLines > funcInfos[j].uncoveredLines
	})

	f := funk.Filter(funcInfos, func(x *funcInfo) bool {
		return x.uncoveredLines > 0
	})

	var fStr [][]string
	fStr = funk.Map(f, func(x *funcInfo) []string {
		return []string{trimString(x.fileName), x.functionName, strconv.FormatInt(x.uncoveredLines, 10)}
	}).([][]string)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Function", "Uncovered Lines"})
	table.AppendBulk(fStr)
	table.Render()
	return
}

func trimString(s string) string {
	if len(s) > 20 {
		return "..." + s[len(s) - 20:]
	}
	return s
}

type funcInfo struct {
	fileName          string
	functionName      string
	functionStartLine int
	uncoveredLines    int64
}

func (f *funcInfo) Print(w io.Writer) {
	_, _ = fmt.Fprintf(w, "%s:%d:\t%s\t%d\n", f.fileName, f.functionStartLine, f.functionName, f.uncoveredLines)
}

// findFuncs parses the file and returns a slice of FuncExtent descriptors.
func findFuncs(name string) ([]*FuncExtent, error) {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}
	visitor := &FuncVisitor{
		fset:    fset,
		name:    name,
		astFile: parsedFile,
	}
	ast.Walk(visitor, visitor.astFile)
	return visitor.funcs, nil
}

// FuncExtent describes a function's extent in the source by file and position.
type FuncExtent struct {
	name      string
	startLine int
	startCol  int
	endLine   int
	endCol    int
}

// FuncVisitor implements the visitor that builds the function position list for a file.
type FuncVisitor struct {
	fset    *token.FileSet
	name    string // Name of file.
	astFile *ast.File
	funcs   []*FuncExtent
}

// Visit implements the ast.Visitor interface.
func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		start := v.fset.Position(n.Pos())
		end := v.fset.Position(n.End())
		fe := &FuncExtent{
			name:      n.Name.Name,
			startLine: start.Line,
			startCol:  start.Column,
			endLine:   end.Line,
			endCol:    end.Column,
		}
		v.funcs = append(v.funcs, fe)
	}
	return v
}

// coverage returns the fraction of the statements in the function that were covered, as a numerator and denominator.
func (f *FuncExtent) coverage(profile *cover.Profile) (num, den int64) {
	// We could avoid making this n^2 overall by doing a single scan and annotating the functions,
	// but the sizes of the data structures is never very large and the scan is almost instantaneous.
	var covered, total int64
	// The blocks are sorted, so we can stop counting as soon as we reach the end of the relevant block.
	for _, b := range profile.Blocks {
		if b.StartLine > f.endLine || (b.StartLine == f.endLine && b.StartCol >= f.endCol) {
			// Past the end of the function.
			break
		}
		if b.EndLine < f.startLine || (b.EndLine == f.startLine && b.EndCol <= f.startCol) {
			// Before the beginning of the function
			continue
		}
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		total = 1 // Avoid zero denominator.
	}
	return covered, total
}

// findFile finds the location of the named file in GOROOT, GOPATH etc.
func findFile(file string) (string, error) {
	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", file, err)
	}
	return filepath.Join(pkg.Dir, file), nil
}
