package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const moduleRoot = "/home/m/wayland"

type TypeDef struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type ParamDef struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FuncDef struct {
	Name     string     `json:"name"`
	RecvType string     `json:"recv_type,omitempty"`
	Params   []ParamDef `json:"params"`
	Results  []string   `json:"results"`
}

type InterfaceDef struct {
	Name    string     `json:"name"`
	Methods []FuncDef  `json:"methods"`
}

type ConstantDef struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PackageAPI struct {
	Platform    string         `json:"platform"`
	PackageName string         `json:"package_name"`
	Types       []TypeDef      `json:"types"`
	Funcs       []FuncDef      `json:"funcs"`
	Interfaces  []InterfaceDef `json:"interfaces"`
	Constants   []ConstantDef  `json:"constants"`
}

type Mismatch struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Report struct {
	Platform    string      `json:"platform"`
	Pkg1        string      `json:"pkg1"`
	Pkg2        string      `json:"pkg2"`
	Mismatches  []Mismatch  `json:"mismatches"`
	Passing     []string    `json:"passing"`
	Summary     string      `json:"summary"`
}

func main() {
	jsonFlag := flag.Bool("json", false, "Output results as JSON")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: apicheck [linux|windows|darwin|js|all]\n")
		fmt.Fprintf(os.Stderr, "  all checks all platforms\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fmt.Fprintf(os.Stderr, "  -json    Output results as JSON\n")
		os.Exit(2)
	}

	targets := args
	if len(targets) == 1 && targets[0] == "all" {
		targets = []string{"linux", "windows", "darwin", "js"}
	}

	module := "github.com/neurlang/wayland"
	pkg1 := module + "/window"
	pkg2 := module + "/windowtrace"

	allReports := []Report{}
	failCount := 0

	for _, target := range targets {
		api1 := loadPackageAPI(target, pkg1)
		api2 := loadPackageAPI(target, pkg2)
		report := compareAPIs(target, pkg1, pkg2, api1, api2)
		allReports = append(allReports, report)
		if len(report.Mismatches) > 0 {
			failCount++
		}
	}

	if *jsonFlag {
		data, _ := json.MarshalIndent(allReports, "", "  ")
		fmt.Println(string(data))
	} else {
		for _, r := range allReports {
			printReport(r)
		}
		fmt.Println("\n=== SUMMARY ===")
		for _, r := range allReports {
			status := "PASS"
			if len(r.Mismatches) > 0 {
				status = "FAIL"
			}
			fmt.Printf("  [%s] %s\n", status, r.Platform)
		}
	}

	if failCount > 0 {
		fmt.Fprintf(os.Stderr, "\n%d platform(s) failed API compatibility check\n", failCount)
		os.Exit(1)
	}
}

func loadPackageAPI(goos, pkgPath string) PackageAPI {
	files := getPackageFiles(goos, pkgPath)
	if len(files) == 0 {
		return PackageAPI{Platform: goos, PackageName: pkgPath}
	}

	fset := token.NewFileSet()
	var allDecls []ast.Decl

	for _, file := range files {
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			continue
		}
		allDecls = append(allDecls, f.Decls...)
	}

	api := PackageAPI{Platform: goos, PackageName: pkgPath}

	for _, decl := range allDecls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.CONST:
				for _, spec := range d.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range vs.Names {
							if !name.IsExported() {
								continue
							}
							c := ConstantDef{Name: name.Name}
							if vs.Type != nil {
								c.Type = exprToString(vs.Type)
							}
							if vs.Values != nil && len(vs.Values) > 0 {
								c.Value = exprToString(vs.Values[0])
							}
							api.Constants = append(api.Constants, c)
						}
					}
				}
			case token.TYPE:
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if ts.Name.IsExported() {
							t := TypeDef{
								Name: ts.Name.Name,
								Kind: exprToString(ts.Type),
							}
							api.Types = append(api.Types, t)
						}
					}
				}
			}
		case *ast.FuncDecl:
			fi := funcDeclToDef(d)
			if fi != nil {
				api.Funcs = append(api.Funcs, *fi)
			}
		}
	}

	// Extract interfaces from type declarations
	for _, decl := range allDecls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
			for _, spec := range gd.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if iface, ok := ts.Type.(*ast.InterfaceType); ok {
						name := ts.Name.Name
						if ts.Name.IsExported() {
							methods := extractInterfaceMethods(iface)
							api.Interfaces = append(api.Interfaces, InterfaceDef{
								Name:    name,
								Methods: methods,
							})
						}
					}
				}
			}
		}
	}

	return api
}

func extractInterfaceMethods(iface *ast.InterfaceType) []FuncDef {
	var methods []FuncDef
	for _, spec := range iface.Methods.List {
		if fd, ok := spec.Type.(*ast.FuncType); ok {
			m := FuncDef{Name: spec.Names[0].Name}
			if fd.Params != nil {
				for _, param := range fd.Params.List {
					for _, name := range param.Names {
						m.Params = append(m.Params, ParamDef{
							Name: name.Name,
							Type: exprToString(param.Type),
						})
					}
				}
			}
			if fd.Results != nil {
				for _, result := range fd.Results.List {
					for range result.Names {
						m.Results = append(m.Results, exprToString(result.Type))
					}
				}
			}
			methods = append(methods, m)
		}
	}
	return methods
}

func funcDeclToDef(fd *ast.FuncDecl) *FuncDef {
	if fd == nil || fd.Name == nil {
		return nil
	}
	if !fd.Name.IsExported() {
		return nil
	}

	fi := &FuncDef{Name: fd.Name.Name}

	if fd.Recv != nil && len(fd.Recv.List) > 0 {
		recv := fd.Recv.List[0]
		recvStr := exprToString(recv.Type)
		// Skip methods on unexported types
		if isUnexported(recvStr) {
			return nil
		}
		fi.RecvType = recvStr
	}

	if fd.Type.Params != nil {
		for _, param := range fd.Type.Params.List {
			for _, name := range param.Names {
				fi.Params = append(fi.Params, ParamDef{
					Name: name.Name,
					Type: exprToString(param.Type),
				})
			}
		}
	}

	if fd.Type.Results != nil {
		for _, result := range fd.Type.Results.List {
			for range result.Names {
				fi.Results = append(fi.Results, exprToString(result.Type))
			}
		}
	}

	return fi
}

func isUnexported(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Strip pointer prefix
	if s[0] == '*' {
		s = s[1:]
	}
	// Check if the first character after * is lowercase
	return s[0] >= 'a' && s[0] <= 'z'
}

func getPackageFiles(goos, pkgPath string) []string {
	var tags string
	switch goos {
	case "darwin":
		tags = "darwin"
	default:
		tags = goos
	}

	cmd := exec.Command("go", "list", "-f", "{{.Dir}}\n{{join .GoFiles \"\\n\"}}", "-tags", tags, pkgPath)
	cmd.Dir = moduleRoot
	cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH=amd64")
	out, err := cmd.Output()
	if err != nil {
		cmd2 := exec.Command("go", "list", "-f", "{{.Dir}}\n{{join .GoFiles \"\\n\"}}", pkgPath)
		cmd2.Dir = moduleRoot
		cmd2.Env = append(os.Environ(), "GOOS="+goos, "GOARCH=amd64")
		out2, err2 := cmd2.Output()
		if err2 != nil {
			return fallbackFiles(pkgPath)
		}
		return parseDirFiles(string(out2))
	}
	return parseDirFiles(string(out))
}

func parseDirFiles(output string) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return nil
	}
	pkgDir := lines[0]
	var files []string
	for _, f := range lines[1:] {
		if f == "" {
			continue
		}
		files = append(files, pkgDir+"/"+f)
	}
	sort.Strings(files)
	return files
}

func fallbackFiles(pkgPath string) []string {
	// Get the actual directory from go list
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", pkgPath)
	cmd.Dir = moduleRoot
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	pkgDir := strings.TrimSpace(string(out))
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return nil
	}
	var result []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".go") {
			result = append(result, pkgDir+"/"+e.Name())
		}
	}
	sort.Strings(result)
	return result
}

func compareAPIs(goos, pkg1, pkg2 string, api1, api2 PackageAPI) Report {
	report := Report{
		Platform: goos,
		Pkg1:     pkg1,
		Pkg2:     pkg2,
	}

	// Build lookup sets
	types1 := make(map[string]bool)
	for _, t := range api1.Types {
		types1[t.Name] = true
	}
	types2 := make(map[string]bool)
	for _, t := range api2.Types {
		types2[t.Name] = true
	}

	funcs1 := make(map[string]bool)
	for _, f := range api1.Funcs {
		funcs1[funcKey(f)] = true
	}
	funcs2 := make(map[string]bool)
	for _, f := range api2.Funcs {
		funcs2[funcKey(f)] = true
	}

	// Build interface maps (name -> sorted method signatures)
	iface1 := buildInterfaceMap(api1)
	iface2 := buildInterfaceMap(api2)

	// Build constants map
	consts1 := make(map[string]bool)
	for _, c := range api1.Constants {
		consts1[c.Name] = true
	}
	consts2 := make(map[string]bool)
	for _, c := range api2.Constants {
		consts2[c.Name] = true
	}

	// Build method maps (recvType.name -> sorted signatures)
	methods1 := buildMethodMap(api1)
	methods2 := buildMethodMap(api2)

	// === CHECK TYPES ===
	for _, t := range api1.Types {
		if !types2[t.Name] {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Severity: "FAIL",
				Message: fmt.Sprintf("Type %q missing in %s", t.Name, pkg2),
			})
		}
	}
	for _, t := range api2.Types {
		if !types1[t.Name] {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Severity: "WARN",
				Message: fmt.Sprintf("Type %q missing in %s", t.Name, pkg1),
			})
		}
	}

	// === CHECK INTERFACES ===
	for name, methods1 := range iface1 {
		if methods2, ok := iface2[name]; ok {
			sort.Strings(methods1)
			sort.Strings(methods2)
			if !equalSlices(methods1, methods2) {
				report.Mismatches = append(report.Mismatches, Mismatch{
					Severity: "FAIL",
					Message: fmt.Sprintf("Interface %q methods differ:\n    %s: %v\n    %s: %v", name, pkg1, methods1, pkg2, methods2),
				})
			} else {
				report.Passing = append(report.Passing, fmt.Sprintf("Interface %q matches", name))
			}
		} else {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Severity: "FAIL",
				Message: fmt.Sprintf("Interface %q missing in %s", name, pkg2),
			})
		}
	}
	for name := range iface2 {
		if _, ok := iface1[name]; !ok {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Severity: "WARN",
				Message: fmt.Sprintf("Interface %q missing in %s", name, pkg1),
			})
		}
	}

	// === CHECK CONSTANTS ===
	for _, c := range api1.Constants {
		if !consts2[c.Name] {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Severity: "FAIL",
				Message: fmt.Sprintf("Constant %q missing in %s", c.Name, pkg2),
			})
		}
	}
	for _, c := range api2.Constants {
		if !consts1[c.Name] {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Severity: "WARN",
				Message: fmt.Sprintf("Constant %q missing in %s", c.Name, pkg1),
			})
		}
	}

	// === CHECK TOP-LEVEL FUNCTIONS ===
	for _, f := range api1.Funcs {
		if f.RecvType == "" {
			found := false
			for _, g := range api2.Funcs {
				if g.Name == f.Name && equalParamTypes(f.Params, g.Params) {
					found = true
					break
				}
			}
			if !found {
				report.Mismatches = append(report.Mismatches, Mismatch{
					Severity: "FAIL",
					Message: fmt.Sprintf("Function %q missing in %s", f.Name, pkg2),
				})
			}
		}
	}
	for _, f := range api2.Funcs {
		if f.RecvType == "" {
			found := false
			for _, g := range api1.Funcs {
				if g.Name == f.Name && equalParamTypes(f.Params, g.Params) {
					found = true
					break
				}
			}
			if !found {
				report.Mismatches = append(report.Mismatches, Mismatch{
					Severity: "WARN",
					Message: fmt.Sprintf("Function %q missing in %s", f.Name, pkg1),
				})
			}
		}
	}

	// === CHECK METHODS ===
	for name := range methods1 {
		if _, ok := methods2[name]; !ok {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Severity: "FAIL",
				Message: fmt.Sprintf("Method %q missing in %s", name, pkg2),
			})
		}
	}
	for name := range methods2 {
		if _, ok := methods1[name]; !ok {
			report.Mismatches = append(report.Mismatches, Mismatch{
				Severity: "WARN",
				Message: fmt.Sprintf("Method %q missing in %s", name, pkg1),
			})
		}
	}
	for name := range methods1 {
		if s1, ok1 := methods1[name]; ok1 {
			if s2, ok2 := methods2[name]; ok2 {
				sort.Strings(s1)
				sort.Strings(s2)
				if !equalSlices(s1, s2) {
					report.Mismatches = append(report.Mismatches, Mismatch{
						Severity: "FAIL",
						Message: fmt.Sprintf("Method %q signature differs:\n    %s: %v\n    %s: %v", name, pkg1, s1, pkg2, s2),
					})
				} else {
					report.Passing = append(report.Passing, fmt.Sprintf("Method %q matches", name))
				}
			}
		}
	}

	// Generate summary
	if len(report.Mismatches) == 0 {
		report.Summary = "API surfaces match"
	} else {
		report.Summary = fmt.Sprintf("%d mismatch(s) found", len(report.Mismatches))
	}

	return report
}

func buildInterfaceMap(api PackageAPI) map[string][]string {
	iface := make(map[string][]string)
	for _, ifaceDef := range api.Interfaces {
		var methods []string
		for _, m := range ifaceDef.Methods {
			methods = append(methods, methodSignature(m))
		}
		sort.Strings(methods)
		iface[ifaceDef.Name] = methods
	}
	return iface
}

func buildMethodMap(api PackageAPI) map[string][]string {
	methods := make(map[string][]string)
	for _, f := range api.Funcs {
		if f.RecvType != "" {
			key := f.RecvType + "." + f.Name
			methods[key] = append(methods[key], funcSignature(f))
		}
	}
	return methods
}

func funcKey(f FuncDef) string {
	return f.Name + "(" + paramsToString(f.Params) + ")"
}

func methodSignature(m FuncDef) string {
	return m.Name + "(" + paramsToString(m.Params) + ")"
}

func funcSignature(f FuncDef) string {
	return f.Name + "(" + paramsToString(f.Params) + ") " + strings.Join(f.Results, ",")
}

func paramsToString(params []ParamDef) string {
	var parts []string
	for _, p := range params {
		parts = append(parts, p.Type)
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func equalParamTypes(a, b []ParamDef) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Type != b[i].Type {
			return false
		}
	}
	return true
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sa := append([]string(nil), a...)
	sb := append([]string(nil), b...)
	sort.Strings(sa)
	sort.Strings(sb)
	for i := range sa {
		if sa[i] != sb[i] {
			return false
		}
	}
	return true
}

func exprToString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.MapType:
		return "map[" + exprToString(e.Key) + "]" + exprToString(e.Value)
	case *ast.ArrayType:
		return "[]" + exprToString(e.Elt)
	case *ast.FuncType:
		return "func"
	case *ast.StructType:
		return "struct{}"
	case *ast.BasicLit:
		return e.Value
	case *ast.UnaryExpr:
		return e.Op.String() + exprToString(e.X)
	case *ast.BinaryExpr:
		return exprToString(e.X) + " " + e.Op.String() + " " + exprToString(e.Y)
	case *ast.ParenExpr:
		return "(" + exprToString(e.X) + ")"
	case *ast.IndexExpr:
		return exprToString(e.X) + "[" + exprToString(e.Index) + "]"
	case *ast.CallExpr:
		return exprToString(e.Fun) + "(...)"
	case *ast.TypeAssertExpr:
		return exprToString(e.X) + ".(" + exprToString(e.Type) + ")"
	case *ast.Ellipsis:
		return "..." + exprToString(e.Elt)
	case *ast.KeyValueExpr:
		return exprToString(e.Key) + ": " + exprToString(e.Value)
	case *ast.IndexListExpr:
		return exprToString(e.X) + "[...]"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func printReport(r Report) {
	fmt.Printf("\n=== Platform: %s ===\n", r.Platform)
	fmt.Printf("  %s\n  %s\n", r.Pkg1, r.Pkg2)

	if len(r.Mismatches) > 0 {
		fmt.Println("\n  MISMATCHES:")
		for _, m := range r.Mismatches {
			fmt.Printf("    [%s] %s\n", m.Severity, m.Message)
		}
	}
	if len(r.Passing) > 0 {
		fmt.Println("\n  PASSING:")
		for _, p := range r.Passing {
			fmt.Printf("    [OK] %s\n", p)
		}
	}

	if len(r.Mismatches) == 0 {
		fmt.Println("\n  ✓ API surfaces match")
	}
}

func init() {
	if _, err := os.Stat(moduleRoot + "/go.mod"); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: go.mod not found at %s\n", moduleRoot+"/go.mod")
		os.Exit(1)
	}
}
