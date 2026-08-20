// Package reflectmethod detects calls to reflect.Value.Method, reflect.Value.MethodByName, reflect.Type.Method, and reflect.Type.MethodByName in a module and its dependencies.
// The compiler flags the enclosing function of these calls with AttrReflectMethod, which makes the linker retain all methods of all reachable types instead of applying dead code elimination.
package reflectmethod

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"slices"
	"sort"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

// Finding describes a call to a reflect method that disables method dead code elimination.
type Finding struct {
	// Module is the module path of the package containing the call.
	Module string
	// Package is the import path of the package containing the call.
	Package string
	// Filename is the absolute path of the file containing the call.
	Filename string
	// Line is the 1-based line of the call.
	Line int
	// Column is the 1-based column of the call.
	Column int
	// Function is the name of the enclosing function, "<anonymous>" for a function literal, or empty for a package-level initializer.
	Function string
	// Method is the fully qualified name of the called method, e.g. "reflect.Value.MethodByName".
	Method string
}

// Detect returns the calls to the reflect methods that disable method dead code elimination (reflect.Value.Method, reflect.Value.MethodByName, reflect.Type.Method, reflect.Type.MethodByName) found in the module containing dir and in its transitive dependencies.
// It fails on compilation errors in the packages it analyzes: syntax errors are reported for every package in the dependency graph, and type errors for the packages it loads for analysis (those that mention a reflect method).
// Type errors in a package that mentions no reflect method are not reported, because such a package is never type-checked.
// Test files are not analyzed.
func Detect(ctx context.Context, dir string) ([]Finding, error) {
	pkgs, err := loadPackages(ctx, dir, []string{"./..."}, packages.NeedName|packages.NeedFiles|packages.NeedSyntax|packages.NeedImports|packages.NeedDeps|packages.NeedModule, "list packages")
	if err != nil {
		return nil, err
	}
	candidates, err := candidatePaths(allPackages(pkgs))
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, nil
	}
	pkgs, err = loadPackages(ctx, dir, candidates, packages.NeedName|packages.NeedFiles|packages.NeedSyntax|packages.NeedTypes|packages.NeedTypesInfo|packages.NeedModule, "load candidate packages")
	if err != nil {
		return nil, err
	}
	return sortFindings(detectInPackages(pkgs)), nil
}

func loadPackages(ctx context.Context, dir string, patterns []string, mode packages.LoadMode, name string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Context: ctx,
		Dir:     dir,
		Mode:    mode,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	loadErrs := make([]error, 0)
	for _, pkg := range pkgs {
		for _, e := range pkg.Errors {
			loadErrs = append(loadErrs, e)
		}
	}
	if len(loadErrs) > 0 {
		return nil, fmt.Errorf("%s: %w", name, errors.Join(loadErrs...))
	}
	return pkgs, nil
}

func allPackages(roots []*packages.Package) []*packages.Package {
	seen := make(map[*packages.Package]bool)
	var out []*packages.Package
	var visit func(pkg *packages.Package)
	visit = func(pkg *packages.Package) {
		if seen[pkg] {
			return
		}
		seen[pkg] = true
		out = append(out, pkg)
		for _, imp := range pkg.Imports {
			visit(imp)
		}
	}
	for _, root := range roots {
		visit(root)
	}
	return out
}

func candidatePaths(pkgs []*packages.Package) ([]string, error) {
	candidates := make(map[string]struct{})
	for _, pkg := range pkgs {
		if pkg.Module == nil {
			continue
		}
		mentions, err := mentionsMethod(pkg)
		if err != nil {
			return nil, fmt.Errorf("check package %s for method lookups: %w", pkg.PkgPath, err)
		}
		if mentions {
			candidates[pkg.PkgPath] = struct{}{}
		}
	}
	paths := make([]string, 0, len(candidates))
	for path := range candidates {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

func mentionsMethod(pkg *packages.Package) (bool, error) {
	for _, filename := range pkg.GoFiles {
		content, err := os.ReadFile(filename)
		if err != nil {
			return false, fmt.Errorf("read file %s: %w", filename, err)
		}
		if bytes.Contains(content, []byte(".Method")) {
			return true, nil
		}
	}
	return false, nil
}

func detectInPackages(pkgs []*packages.Package) []Finding {
	findings := make([]Finding, 0)
	for _, pkg := range pkgs {
		insp := inspector.New(pkg.Syntax)
		insp.WithStack([]ast.Node{(*ast.CallExpr)(nil), (*ast.FuncDecl)(nil), (*ast.FuncLit)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
			if !push {
				return true
			}
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := unwrapParen(call.Fun).(*ast.SelectorExpr)
			if !ok || !isReflectMethodLookup(sel) {
				return true
			}
			recv, ok := reflectMethodReceiver(pkg.TypesInfo, sel)
			if !ok {
				return true
			}
			findings = append(findings, newFinding(pkg, sel, recv, stack))
			return true
		})
	}
	return findings
}

func isReflectMethodLookup(sel *ast.SelectorExpr) bool {
	return sel.Sel.Name == "Method" || sel.Sel.Name == "MethodByName"
}

// unwrapParen returns e without any enclosing parentheses.
func unwrapParen(e ast.Expr) ast.Expr {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.X
	}
}

func newFinding(pkg *packages.Package, sel *ast.SelectorExpr, recv *types.Named, stack []ast.Node) Finding {
	pos := pkg.Fset.Position(sel.Sel.Pos())
	return Finding{
		Module:   modulePath(pkg),
		Package:  pkg.PkgPath,
		Filename: pos.Filename,
		Line:     pos.Line,
		Column:   pos.Column,
		Function: enclosingFunc(stack),
		Method:   methodName(recv, sel),
	}
}

func reflectMethodReceiver(info *types.Info, sel *ast.SelectorExpr) (*types.Named, bool) {
	fn, ok := info.ObjectOf(sel.Sel).(*types.Func)
	if !ok {
		return nil, false
	}
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Recv() == nil {
		return nil, false
	}
	recv, ok := namedType(sig.Recv().Type())
	if !ok {
		return nil, false
	}
	obj := recv.Obj()
	if obj.Pkg() == nil || obj.Pkg().Path() != "reflect" {
		return nil, false
	}
	if obj.Name() != "Value" && obj.Name() != "Type" {
		return nil, false
	}
	return recv, true
}

func namedType(t types.Type) (*types.Named, bool) {
	if ptr, isPtr := t.(*types.Pointer); isPtr {
		t = ptr.Elem()
	}
	named, ok := t.(*types.Named)
	return named, ok
}

func methodName(recv *types.Named, sel *ast.SelectorExpr) string {
	obj := recv.Obj()
	if obj.Pkg() == nil {
		return sel.Sel.Name
	}
	return obj.Pkg().Name() + "." + obj.Name() + "." + sel.Sel.Name
}

func enclosingFunc(stack []ast.Node) string {
	for _, v := range slices.Backward(stack) {
		switch n := v.(type) {
		case *ast.FuncLit:
			return "<anonymous>"
		case *ast.FuncDecl:
			return n.Name.Name
		}
	}
	return ""
}

func modulePath(pkg *packages.Package) string {
	if pkg.Module == nil {
		return ""
	}
	return pkg.Module.Path
}

func sortFindings(findings []Finding) []Finding {
	slices.SortFunc(findings, func(a, b Finding) int {
		return cmp.Or(
			cmp.Compare(a.Filename, b.Filename),
			cmp.Compare(a.Line, b.Line),
			cmp.Compare(a.Column, b.Column),
		)
	})
	return findings
}
