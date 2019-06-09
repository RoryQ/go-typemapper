package generator

import (
	"bytes"
	"flag"
	"fmt"
	"go/types"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

var update = flag.Bool("update", false, "update generated files")

func loadSSAPackage(t *testing.T, path string) *ssa.Package {
	t.Helper()

	cfg := &packages.Config{
		Dir:        path,
		Mode:       packages.LoadSyntax,
		BuildFlags: []string{"-tags", BuildTag},
		Tests:      false,
	}
	pkgs, err := packages.Load(cfg, "./")
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("error: %d packages found", len(pkgs))
	}
	pkg := pkgs[0]

	prog := ssa.NewProgram(pkg.Fset, ssa.BuilderMode(0))

	// Create SSA packages for all imports.
	// Order is not significant.
	created := make(map[*types.Package]bool)
	var createAll func(pkgs []*types.Package)
	createAll = func(pkgs []*types.Package) {
		for _, p := range pkgs {
			if !created[p] {
				created[p] = true
				prog.CreatePackage(p, nil, nil, true)
				createAll(p.Imports())
			}
		}
	}
	createAll(pkg.Types.Imports())

	// Create and build the primary package.
	ssapkg := prog.CreatePackage(pkg.Types, pkg.Syntax, pkg.TypesInfo, false)
	ssapkg.Build()

	return ssapkg
}

func TestExamples(t *testing.T) {
	examplesPath, err := filepath.Abs("../examples")
	if err != nil {
		t.Fatal(err)
	}

	examples, err := ioutil.ReadDir(examplesPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, ex := range examples {
		c := ex
		t.Run(c.Name(), func(t *testing.T) {
			assert := require.New(t)

			if !c.IsDir() {
				t.Fatalf("%q is not a directory", c.Name())
			}
			pkgPath := filepath.Join(examplesPath, c.Name())
			ssapkg := loadSSAPackage(t, pkgPath)

			g := NewGenerator(
				ssapkg,
				"// Code generated by \"typemapper \"; DO NOT EDIT.\n",
				fmt.Sprintf("// +build !%s", BuildTag),
			)
			err = g.GenerateMappings()
			assert.NoError(err)

			for _, fileName := range g.AllFiles() {
				outFile := filepath.Join(pkgPath, fileName)
				expectedBytes, err := ioutil.ReadFile(outFile)
				assert.NoError(err)
				expected := strings.ReplaceAll(string(expectedBytes), "\r\n", "\n")

				actualBuf := &bytes.Buffer{}
				err = g.Render(fileName, actualBuf)
				assert.NoError(err)

				if *update {
					err = ioutil.WriteFile(outFile, actualBuf.Bytes(), 0777)
					assert.NoError(err)
				} else {
					actual := strings.ReplaceAll(actualBuf.String(), "\r\n", "\n")
					assert.Equal(expected, actual)
				}
			}
		})
	}
}
