package reflectmethod

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"sort"
	"testing"

	"github.com/pierrre/assert"
	"golang.org/x/tools/go/packages"
)

func setGOWORKOff(t *testing.T) {
	t.Helper()
	// GOWORK=off isolates the testdata modules from any go.work: with a workspace active, packages.Load silently resolves nothing and the test reports zero findings.
	t.Setenv("GOWORK", "off")
}

func TestDetect(t *testing.T) {
	setGOWORKOff(t)
	ctx := t.Context()
	findings, err := Detect(ctx, filepath.Join("testdata", "main"))
	assert.NoError(t, err)
	got := make([]string, 0, len(findings))
	for _, f := range findings {
		got = append(got, fmt.Sprintf("%s %d:%d %s %s %s %s", filepath.Base(f.Filename), f.Line, f.Column, f.Module, f.Package, f.Function, f.Method))
	}
	sort.Strings(got)
	assert.SliceEqual(t, []string{
		"dep.go 10:37 example.com/dep example.com/dep First reflect.Value.Method",
		"dep.go 6:37 example.com/dep example.com/dep Call reflect.Value.MethodByName",
		"main.go 33:46 example.com/main example.com/main  reflect.Value.Method",
		"main.go 41:4 example.com/main example.com/main main reflect.Value.Method",
		"main.go 42:4 example.com/main example.com/main main reflect.Value.MethodByName",
		"main.go 44:4 example.com/main example.com/main main reflect.Type.Method",
		"main.go 45:4 example.com/main example.com/main main reflect.Type.MethodByName",
		"main.go 47:8 example.com/main example.com/main main reflect.Value.Method",
		"main.go 48:16 example.com/main example.com/main main reflect.Value.MethodByName",
		"main.go 52:29 example.com/main example.com/main <anonymous> reflect.Value.Method",
		"main.go 54:17 example.com/main example.com/main main reflect.Value.Method",
		"other.go 6:37 example.com/main example.com/main other reflect.Value.Method",
	}, got)
}

func TestDetectClean(t *testing.T) {
	setGOWORKOff(t)
	ctx := t.Context()
	findings, err := Detect(ctx, filepath.Join("testdata", "clean"))
	assert.NoError(t, err)
	assert.SliceEmpty(t, findings)
}

func TestDetectBroken(t *testing.T) {
	setGOWORKOff(t)
	ctx := t.Context()
	_, err := Detect(ctx, filepath.Join("testdata", "broken"))
	assert.ErrorContains(t, err, "main.go:3")
}

func TestDetectBrokenType(t *testing.T) {
	setGOWORKOff(t)
	ctx := t.Context()
	_, err := Detect(ctx, filepath.Join("testdata", "broken_type"))
	assert.ErrorContains(t, err, "cannot use")
}

func TestLoadPackagesError(t *testing.T) {
	ctx := t.Context()
	_, err := loadPackages(ctx, filepath.Join(t.TempDir(), "nonexistent"), []string{"./..."}, packages.NeedName, "list packages")
	assert.Error(t, err)
}

func TestMentionsMethodError(t *testing.T) {
	pkg := &packages.Package{GoFiles: []string{filepath.Join(t.TempDir(), "missing.go")}}
	_, err := mentionsMethod(pkg)
	assert.Error(t, err)
}

func TestCandidatePathsError(t *testing.T) {
	pkg := &packages.Package{
		PkgPath: "example.com/p",
		Module:  &packages.Module{Path: "example.com/m"},
		GoFiles: []string{filepath.Join(t.TempDir(), "missing.go")},
	}
	_, err := candidatePaths([]*packages.Package{pkg})
	assert.Error(t, err)
}

func TestNamedType(t *testing.T) {
	named := types.NewNamed(types.NewTypeName(0, types.NewPackage("p", "p"), "T", nil), types.Typ[types.Int], nil)
	got, ok := namedType(named)
	assert.True(t, ok)
	assert.True(t, got == named)
	got, ok = namedType(types.NewPointer(named))
	assert.True(t, ok)
	assert.True(t, got == named)
	_, ok = namedType(types.Typ[types.Int])
	assert.False(t, ok)
	_, ok = namedType(types.NewPointer(types.Typ[types.Int]))
	assert.False(t, ok)
}

func TestMethodName(t *testing.T) {
	// error is in the universe scope, so its Obj().Pkg() is nil.
	errType, ok := assert.Type[*types.Named](t, types.Universe.Lookup("error").Type())
	assert.True(t, ok)
	sel := &ast.SelectorExpr{Sel: ast.NewIdent("Method")}
	assert.Equal(t, "Method", methodName(errType, sel))
}

func TestModulePath(t *testing.T) {
	assert.Equal(t, "", modulePath(&packages.Package{}))
	assert.Equal(t, "example.com/x", modulePath(&packages.Package{Module: &packages.Module{Path: "example.com/x"}}))
}

func TestUnwrapParen(t *testing.T) {
	sel := &ast.SelectorExpr{}
	assert.True(t, unwrapParen(sel) == sel)
	assert.True(t, unwrapParen(&ast.ParenExpr{X: sel}) == sel)
	assert.True(t, unwrapParen(&ast.ParenExpr{X: &ast.ParenExpr{X: sel}}) == sel)
}

func TestReflectMethodReceiver(t *testing.T) {
	sel := &ast.SelectorExpr{Sel: ast.NewIdent("Method")}
	reflectPkg := types.NewPackage("reflect", "reflect")

	// the selector resolves to a non-function (e.g. a struct field).
	info := &types.Info{Uses: map[*ast.Ident]types.Object{sel.Sel: types.NewVar(0, reflectPkg, "Method", types.Typ[types.Int])}}
	_, ok := reflectMethodReceiver(info, sel)
	assert.False(t, ok)

	// the receiver is not a named type.
	recv := types.NewVar(0, nil, "r", types.Typ[types.Int])
	fn := types.NewFunc(0, reflectPkg, "Method", types.NewSignatureType(recv, nil, nil, nil, nil, false))
	info = &types.Info{Uses: map[*ast.Ident]types.Object{sel.Sel: fn}}
	_, ok = reflectMethodReceiver(info, sel)
	assert.False(t, ok)

	// the receiver is a named type that is neither reflect.Value nor reflect.Type.
	other := types.NewNamed(types.NewTypeName(0, reflectPkg, "Other", nil), types.Typ[types.Int], nil)
	recv = types.NewVar(0, nil, "r", other)
	fn = types.NewFunc(0, reflectPkg, "Method", types.NewSignatureType(recv, nil, nil, nil, nil, false))
	info = &types.Info{Uses: map[*ast.Ident]types.Object{sel.Sel: fn}}
	_, ok = reflectMethodReceiver(info, sel)
	assert.False(t, ok)

	// the receiver is reflect.Value.
	value := types.NewNamed(types.NewTypeName(0, reflectPkg, "Value", nil), types.Typ[types.Int], nil)
	recv = types.NewVar(0, nil, "r", value)
	fn = types.NewFunc(0, reflectPkg, "Method", types.NewSignatureType(recv, nil, nil, nil, nil, false))
	info = &types.Info{Uses: map[*ast.Ident]types.Object{sel.Sel: fn}}
	got, ok := reflectMethodReceiver(info, sel)
	assert.True(t, ok)
	assert.True(t, got == value)
}
