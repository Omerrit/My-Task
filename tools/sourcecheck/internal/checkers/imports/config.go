package imports

import (
	"bufio"
	"bytes"
	"fmt"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/tools/sourcecheck/internal/filter"
	"gerrit-share.lan/go/utils/maps"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"unicode"
	"unicode/utf8"
)

const importFileName = "imports.sc"

const (
	nonPolicy   = ""
	allowPolicy = "Allow"
	denyPolicy  = "Deny"
)

type config struct {
	RootInternalDir  string
	PackageName      string
	DirName          string
	PatternsForAllow []string
	PatternsForDeny  []string
	Variables        maps.String
}

var (
	varMatcher      = regexp.MustCompile(`^(\w+)*[ \t]*=[ \t]*([@\?\*\./\w]+)*`)
	varNameMatcher  = regexp.MustCompile(`^@(\w+)`)
	importValidator = regexp.MustCompile(`^\w+[\w\.\-]*\w+(/\w+[\w\.\-]*\w+)*`)
)

func newConfig(dirName string, packageName string, parentDirSettings dirSettings) (*config, error) {
	config := config{
		RootInternalDir: parentDirSettings.RootInternalDir,
		PackageName:     packageName,
		DirName:         dirName,
		Variables:       parentDirSettings.Variables.Clone(),
	}
	if filepath.Base(dirName) == "internal" {
		config.RootInternalDir = filepath.Dir(dirName)
	}
	err := config.loadFromImportFile(dirName)
	if err != nil {
		return nil, err
	}
	err = config.substituteTemplate()
	if err != nil {
		return nil, err
	}
	err = config.validateAllowedImports()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (conf *config) MakePackageSettings() (packageSettings, error) {
	var (
		packageSettings packageSettings
		err             error
	)
	for _, allowedImport := range conf.PatternsForAllow {
		packageSettings.AllowedImports.Add(allowedImport)
	}
	packageSettings.Deny, err = filter.NewFilterSet(conf.PatternsForDeny)
	if err != nil {
		return packageSettings, err
	}
	packageSettings.DirName = conf.DirName
	return packageSettings, nil
}

func (conf *config) MakeDirSettings() dirSettings {
	return dirSettings{
		RootInternalDir: conf.RootInternalDir,
		ImportPath:      conf.PackageName,
		Variables:       conf.Variables,
	}
}

func (conf *config) substituteTemplate() error {
	var err error
	conf.PatternsForAllow, err = substituteTemplate(conf.PatternsForAllow, conf.PackageName, conf.Variables)
	if err != nil {
		return err
	}
	conf.PatternsForDeny, err = substituteTemplate(conf.PatternsForDeny, conf.PackageName, conf.Variables)
	if err != nil {
		return err
	}
	return nil
}

func substituteTemplate(patterns []string, packageName string, variables maps.String) ([]string, error) {
	var err error
	patterns, err = substituteVariables(patterns, variables)
	if err != nil {
		return nil, err
	}
	patterns = substituteCurrentPackageName(patterns, packageName)
	return patterns, nil
}

func substituteVariables(patterns []string, variables maps.String) ([]string, error) {
	for i, pattern := range patterns {
		patterns[i] = substituteVarNameForValue(pattern, variables)
	}
	return patterns, nil
}

func substituteCurrentPackageName(patterns []string, packageName string) []string {
	for i, pattern := range patterns {
		if strings.HasPrefix(pattern, "./") {
			patterns[i] = packageName + pattern[1:]
		}
	}
	return patterns
}

func (conf *config) validateAllowedImports() error {
	var invalidImports []string
	for _, allowedImport := range conf.PatternsForAllow {
		if importValidator.FindString(allowedImport) != allowedImport {
			invalidImports = append(invalidImports, allowedImport)
		}
	}
	if len(invalidImports) > 0 {
		messageBuilder := strings.Builder{}
		messageBuilder.WriteString(fmt.Sprintf("in file %s: invalid allowed imports:\n",
			conf.DirName+"/"+importFileName))
		for _, invalidImport := range invalidImports {
			messageBuilder.WriteString(invalidImport + "\n")
		}
		return fmt.Errorf(messageBuilder.String())
	}
	return nil
}

func (conf *config) loadFromImportFile(dir string) error {
	fileName := filepath.Join(dir, importFileName)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		// file not found
		if errors.Is(err, syscall.ENOENT) {
			return nil
		}
		return err
	}
	if err := conf.parseData(data); err != nil {
		return fmt.Errorf("in file %s: %w", conf.DirName+"/"+importFileName, err)
	}
	return nil
}

func isVariable(s string) (name string, value string, ok bool) {
	matchList := varMatcher.FindStringSubmatch(s)
	if matchList == nil {
		return "", "", false
	}
	return matchList[1], matchList[2], true
}

func substituteVarNameForValue(pattern string, variables maps.String) string {
	if matchList := varNameMatcher.FindStringSubmatchIndex(pattern); matchList != nil {
		varBegin, varEnd := matchList[2], matchList[3]
		pattern = variables[pattern[varBegin:varEnd]] + pattern[varEnd:]
	}
	return pattern
}

func (conf *config) parseData(data []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(scanWords)
	policy := nonPolicy
	blockIsStarted := false
	line := 1
	for scanner.Scan() {
		word := scanner.Text()
		switch word {
		case "\n":
			line++
			blockIsStarted = false
		case allowPolicy, denyPolicy:
			if blockIsStarted {
				return fmt.Errorf("in line %v: could not use %s and %s in same line", line, allowPolicy, denyPolicy)
			}
			policy = word
			blockIsStarted = true
		default:
			if name, value, ok := isVariable(word); ok {
				if blockIsStarted {
					return fmt.Errorf("in line %v: could not define variable and policy in same line", line)
				}
				value = substituteVarNameForValue(value, conf.Variables)
				conf.Variables.Add(name, value)
				policy = nonPolicy
				continue
			}
			switch policy {
			case allowPolicy:
				conf.PatternsForAllow = append(conf.PatternsForAllow, word)
			case denyPolicy:
				conf.PatternsForDeny = append(conf.PatternsForDeny, word)
			default:
				return fmt.Errorf("in line %v: unknown keyword, expect %q or %q, but found %q",
					line, allowPolicy, denyPolicy, word)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("in line %v: %w", line, err)
	}
	return nil
}

func skipSpacesAndComment(data []byte) ([]byte, int) {
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) {
			if r != '#' {
				return data[start:], start
			}
			if lenComment := bytes.IndexRune(data[start:], '\n'); lenComment > 0 {
				return data[start+lenComment:], start + lenComment
			}
			// all data is comment
			return nil, len(data)
		} else {
			if r == '\n' {
				return data[start:], start
			}
		}
	}
	return nil, start
}

func isKeyWord(s string) bool {
	return s == allowPolicy || s == denyPolicy
}

func scanVariable(data []byte) (advance int, token []byte, err error) {
	if matchList := varMatcher.FindSubmatchIndex(data); matchList != nil {
		if matchList[2] < 0 { // variable name did not find
			return 0, nil, fmt.Errorf("except varible name, but assign found")
		}
		token = data[matchList[2]:matchList[3]]
		if isKeyWord(string(token)) {
			return 0, nil, fmt.Errorf("varible name could not be keyword %q", string(token))
		}
		token = append(token, '=')
		if matchList[4] > 0 { // if found value
			token = append(token, data[matchList[4]:matchList[5]]...)
		}
		return matchList[1], token, nil
	}
	// not found
	return 0, nil, nil
}

func scanWords(data []byte, _ bool) (advance int, token []byte, err error) {
	//fmt.Println("Data befor skipping:", string(data))
	data, advance = skipSpacesAndComment(data)
	//fmt.Println("Data after skipping:", string(data),"advance:", advance)
	scanAdvance, token, err := scanVariable(data)
	if err != nil {
		return advance + scanAdvance, nil, err
	}
	if token != nil {
		return advance + scanAdvance, token, nil
	}
	scanAdvance, token, err = scanToken(data)
	if err != nil {
		return advance + scanAdvance, nil, err
	}
	return advance + scanAdvance, token, nil
}

func scanToken(data []byte) (advance int, token []byte, err error) {
	r, width := utf8.DecodeRune(data)
	if unicode.IsSpace(r) {
		if r == '\n' {
			return width, data[:width], nil
		} else {
			return 0, nil, fmt.Errorf("the string started with a space character")
		}
	}
	for i := width; i < len(data); i += width {
		r, width = utf8.DecodeRune(data[i:])
		if unicode.IsSpace(r) {
			return i, data[:i], nil
		}
	}
	return len(data), data, nil
}
