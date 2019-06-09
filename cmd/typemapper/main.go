package main

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"

	"github.com/paultyng/go-typemapper/generator"
)

func main() {
	err := mainErr()
	if err != nil {
		log.Fatal(err)
	}
}

func mainErr() error {
	cfg := &packages.Config{
		Mode:       packages.LoadSyntax,
		BuildFlags: []string{"-tags", generator.BuildTag},
		Tests:      false,
	}
	pkgs, err := packages.Load(cfg, "./")
	if err != nil {
		return err
	}
	if len(pkgs) != 1 {
		return errors.Errorf("error: %d packages found", len(pkgs))
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

	//ssapkg.WriteTo(os.Stdout)

	g := generator.NewGenerator(
		ssapkg,
		// code generation comment
		fmt.Sprintf("// Code generated by \"typemapper %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " ")),
		// inverse of typemapper build tags
		fmt.Sprintf("// +build !%s", generator.BuildTag),
	)

	err = g.GenerateMappings()
	if err != nil {
		return err
	}

	for _, fileName := range g.AllFiles() {
		fileName := fileName
		err = func() error {
			outFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
			if err != nil {
				return err
			}
			defer outFile.Close()

			err = g.Render(fileName, outFile)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
