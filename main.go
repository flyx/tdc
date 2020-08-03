package main

import (
	"os"
	"path/filepath"

	"github.com/flyx/askew/packages"

	"github.com/flyx/askew/data"

	"github.com/pborman/getopt/v2"
)

func main() {
	output := getopt.StringLong(
		"outputDir", 'o', ".", "output directory for index.html")
	initPath := getopt.StringLong(
		"initPath", 'i', "init.go", "path where the Go initialization code should be written to")
	initPkg := getopt.StringLong(
		"initPkg", 'p', "main", "name of the package for the Go initialization code")
	skeletonPath := getopt.StringLong(
		"skeletonPath", 's', "", "path to the skeleton HTML file")
	getopt.Parse()
	var err error
	outputDirPath, err := filepath.Abs(*output)
	if err != nil {
		panic(err)
	}
	if *skeletonPath == "" {
		// will be resolved relative to given working path
		*skeletonPath = "skeleton.html"
	} else {
		if *skeletonPath, err = filepath.Abs(*skeletonPath); err != nil {
			panic(err)
		}
	}

	args := getopt.Args()
	if len(args) == 1 {
		if err := os.Chdir(args[0]); err != nil {
			os.Stdout.WriteString("[error] cannot process directory: " + err.Error() + "\n")
			os.Exit(1)
		}
	} else if len(args) > 0 {
		os.Stdout.WriteString("[error] unexpected arguments:\n")
		for i := 1; i < len(args); i++ {
			os.Stdout.WriteString("[error]   " + args[i] + "\n")
		}
		os.Exit(1)
	}

	info, err := os.Stat(*output)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(*output, os.ModePerm)
			if err != nil {
				panic("unable to create output directory " + *output)
			}
		} else {
			panic("unable to access output directory " + *output)
		}
	} else if !info.IsDir() {
		panic("output path is not a directory: " + *output)
	}

	base, err := packages.Discover()
	if err != nil {
		os.Stdout.WriteString("[error] " + err.Error() + "\n")
		os.Exit(1)
	}
	order, err := packages.Sort(base.Packages)
	if err != nil {
		os.Stdout.WriteString("[error] " + err.Error() + "\n")
		os.Exit(1)
	}

	var p processor
	p.init(base)
	for _, path := range order {
		if err := p.processMacros(path); err != nil {
			os.Stdout.WriteString("[error] " + err.Error() + "\n")
			os.Exit(1)
		}
	}
	for _, path := range order {
		if err := p.processComponents(path); err != nil {
			os.Stdout.WriteString("[error] " + err.Error() + "\n")
			os.Exit(1)
		}
	}

	var s *data.Skeleton
	if s, err = readSkeleton(&p.syms, *skeletonPath); err != nil {
		os.Stdout.WriteString("[error] " + err.Error() + "\n")
		os.Exit(1)
	}

	os.Stdout.WriteString("[info] generating code\n")
	p.dump(s, outputDirPath, *initPath, *initPkg)
}
