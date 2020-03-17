package main

import (
	"flag"
	"fmt"
	"gerrit-share.lan/go/tools/sourcecheck/internal/checkers"
	"gerrit-share.lan/go/tools/sourcecheck/internal/checkers/format"
	"gerrit-share.lan/go/tools/sourcecheck/internal/checkers/imports"
	"gerrit-share.lan/go/tools/sourcecheck/internal/checkers/variables"
	"gerrit-share.lan/go/utils/flags"
	"gerrit-share.lan/go/utils/maps"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type parameters struct {
	export        bool
	format        bool
	badImport     bool
	useCurrentDir bool
	useTopCheck   bool
}

func (p *parameters) registerFlags() {
	flags.BoolFlag(&p.export, "e", "check exported variables")
	flags.BoolFlag(&p.badImport, "i", "check bad imports")
	flags.BoolFlag(&p.format, "f", "check if the file is formatted (using 'gofmt')")
	flags.BoolFlag(&p.useCurrentDir, "c", "start check from current directory (default from repository)")
	flags.BoolFlag(&p.useTopCheck, "t", "check only this directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [dir]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
}

func (p *parameters) useAllCheckers() bool {
	return !(p.badImport || p.export || p.format)
}

func findInputtedDir() string {
	if args := flag.Args(); len(args) > 0 {
		return args[0]
	}
	return "."
}

func findStartDir(params parameters) (string, error) {
	dir := findInputtedDir()
	startDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	if params.useCurrentDir {
		return startDir, nil
	}
	reposDir, err := findGitRepos(startDir)
	if err != nil {
		return "", err
	}
	if reposDir != "" {
		startDir = reposDir
	}
	return startDir, nil
}

func findSubDir(dir string, buffer []string) ([]string, error) {
	listFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range listFiles {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			buffer = append(buffer, filepath.Join(dir, file.Name()))
		}
	}
	return buffer, nil
}

func findAllSubDirs(startDir string) ([]string, error) {
	listDirs := []string{startDir}
	subDirs, err := findSubDir(startDir, nil)
	if err != nil {
		return nil, err
	}
	listDirs = append(listDirs, subDirs...)
	for len(subDirs) > 0 {
		oldSubDirs := make([]string, len(subDirs))
		copy(oldSubDirs, subDirs)
		subDirs = nil
		for _, dir := range oldSubDirs {
			subDirs, err = findSubDir(dir, subDirs)
			if err != nil {
				return nil, err
			}
		}
		listDirs = append(listDirs, subDirs...)
	}
	return listDirs, nil
}

func findGoFiles(dir string) ([]string, error) {
	listFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range listFiles {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			fileNames = append(fileNames, filepath.Join(dir, file.Name()))
		}
	}
	return fileNames, nil
}

func findPackageImportPath(dir string) (string, error) {
	pkg, err := build.ImportDir(dir, build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkg.ImportPath, nil
}

func makeCheckers(params parameters) []checkers.Checker {
	if params.useAllCheckers() {
		return []checkers.Checker{
			imports.NewImportChecker(),
			variables.NewExportVarsChecker(),
			format.NewFormatChecker(),
		}
	}
	var checkers []checkers.Checker
	if params.badImport {
		checkers = append(checkers, imports.NewImportChecker())
	}
	if params.export {
		checkers = append(checkers, variables.NewExportVarsChecker())
	}
	if params.format {
		checkers = append(checkers, format.NewFormatChecker())
	}
	return checkers
}

func checkDir(dir string, packageName string, fileCheckers []checkers.Checker) (info, error) {
	fileNames, err := findGoFiles(dir)
	if err != nil {
		return info{}, err
	}
	var output checkers.Messages
	for _, fileName := range fileNames {
		for _, checker := range fileCheckers {
			output.AppendMessages(checker.Check(fileName))
		}
	}
	return info{
		messages:          output,
		packageImportPath: packageName,
	}, nil
}

func setup(dirs []string, fileCheckers []checkers.Checker) (maps.String, error) {
	var dirPackageMap maps.String
	for _, dir := range dirs {
		packageImportPath, err := findPackageImportPath(dir)
		if err != nil {
			return nil, err
		}
		for _, checker := range fileCheckers {
			err := checker.Setup(dir, packageImportPath)
			if err != nil {
				return nil, err
			}
		}
		dirPackageMap.Add(dir, packageImportPath)
	}
	return dirPackageMap, nil
}

func check(dirs []string, dirPackageMap maps.String, fileCheckers []checkers.Checker) (information, error) {
	var allInfo information
	for _, dir := range dirs {
		info, err := checkDir(dir, dirPackageMap[dir], fileCheckers)
		if err != nil {
			return nil, err
		}
		if !info.IsEmpty() {
			allInfo = append(allInfo, info)
		}
	}
	return allInfo, nil
}

func main() {
	params := parameters{}
	params.registerFlags()
	flag.Parse()
	// find repository directory if it is exist else use inputed directory
	startDir, err := findStartDir(params)
	if err != nil {
		log.Fatalln(err)
	}
	var dirs []string
	if params.useTopCheck {
		dirs = append(dirs, startDir)
	} else {
		dirs, err = findAllSubDirs(startDir)
		if err != nil {
			log.Fatalln(err)
		}
	}
	//for debuging
	//for _, dir := range dirs {
	//	fmt.Println(dir)
	//}
	//return
	checkers := makeCheckers(params)
	dirPackageMap, err := setup(dirs, checkers)
	//_, err = setup(dirs, checkers)
	if err != nil {
		log.Fatalln(err)
	}
	information, err := check(dirs, dirPackageMap, checkers)
	if err != nil {
		log.Fatalln(err)
	}
	if !information.IsEmpty() {
		fmt.Println(information)
		os.Exit(1)
	}
}
