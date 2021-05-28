package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"golang.org/x/tools/cover"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var trim bool

	app := &cli.App{
		Name:  "go-coverage",
		Usage: "identify complex untested functions",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:  "line-filter",
				Value: 0,
				Usage: "functions with untested lines lower than this will be filtered out",
			},
			&cli.BoolFlag{
				Name:        "trim",
				Aliases:     []string{"t"},
				Value:       false,
				Usage:       "trim file name",
				Destination: &trim,
			},
			&cli.StringFlag{
				Name:  "format",
				Value: "table",
				Usage: "display format",
			},
			&cli.StringFlag{
				Name:  "exclude",
				Value: "",
				Usage: "regex of the file to exclude",
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "coverage file",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			profiles, err := cover.ParseProfiles(c.String("file"))
			if err != nil {
				return err
			}

			funcInfos, total, covered, err := getFunctionInfos(profiles)
			if err != nil {
				return err
			}

			sort.Slice(funcInfos, func(i, j int) bool {
				return funcInfos[i].uncoveredLines > funcInfos[j].uncoveredLines
			})

			f := funk.Filter(funcInfos, func(x *funcInfo) bool {
				return x.uncoveredLines > c.Int64("line-filter")
			}).([]*funcInfo)

			exc := c.String("exclude")

			if exc != "" {
				r, regexErr := regexp.Compile(exc)
				if regexErr != nil {
					return regexErr
				}

				f = funk.Filter(f, func(x *funcInfo) bool {
					return !r.Match([]byte(x.fileName))
				}).([]*funcInfo)
			}

			if c.String("format") == "table" {
				printTable(f, trim, covered, total)
			} else {
				printBat(f, trim, covered, total)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getTrimmedFileName(x *funcInfo, trim bool) string {
	fn := x.fileName

	if trim {
		fn = trimString(fn, 20)
	}
	return fn
}

func printBat(f []*funcInfo, trim bool, covered int64, total int64) {
	var fStr [][]string
	tc := float64(covered) / float64(total) * 100
	fStr = funk.Map(f, func(x *funcInfo) []string {
		fn := getTrimmedFileName(x, trim)
		return []string{
			fn,
			x.functionName,
			strconv.Itoa(x.functionStartLine),
			strconv.Itoa(x.functionEndLine),
			strconv.FormatInt(x.uncoveredLines, 10),
			fmt.Sprintf("%.1f", (float64(covered+x.uncoveredLines)/float64(total)*100)-tc)}
	}).([][]string)

	for _, v := range fStr {
		fmt.Println(strings.Join(v, " "))
	}
}

func printTable(f []*funcInfo, trim bool, covered int64, total int64) {
	var fStr [][]string
	tc := float64(covered) / float64(total) * 100
	fStr = funk.Map(f, func(x *funcInfo) []string {
		fn := getTrimmedFileName(x, trim)
		return []string{
			fn,
			x.functionName,
			strconv.FormatInt(x.uncoveredLines, 10),
			fmt.Sprintf("%.1f", (float64(covered+x.uncoveredLines)/float64(total)*100)-tc)}
	}).([][]string)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Function", "Uncovered Lines", "Impact"})
	table.AppendBulk(fStr)
	table.Render()
}

func getFunctionInfos(profiles []*cover.Profile) ([]*funcInfo, int64, int64, error) {
	var total, covered int64
	var funcInfos []*funcInfo
	for _, profile := range profiles {
		fn := profile.FileName
		file, err := findFile(fn)
		if err != nil {
			return nil, 0, 0, err
		}
		funcs, err := findFuncs(file)
		if err != nil {
			return nil, 0, 0, err
		}
		// Now match up functions and profile blocks.
		for _, f := range funcs {
			c, t := f.coverage(profile)
			funcInfos = append(funcInfos,
				&funcInfo{fileName: file,
					functionName:      f.name,
					functionStartLine: f.startLine,
					functionEndLine:   f.endLine,
					uncoveredLines:    t - c})
			total += t
			covered += c
		}
	}
	return funcInfos, total, covered, nil
}

func trimString(s string, i int) string {
	if len(s) > i {
		return "..." + s[len(s)-i:]
	}
	return s
}

type funcInfo struct {
	fileName          string
	functionName      string
	functionStartLine int
	functionEndLine   int
	uncoveredLines    int64
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
