package os2

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/analog-substance/fileutil"
	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
	"github.com/bmatcuk/doublestar/v4"
)

type module struct {
	getCompiled func() (*tengo.Compiled, context.Context)
}

func Module(getCompiled func() (*tengo.Compiled, context.Context)) map[string]tengo.Object {
	m := &module{
		getCompiled: getCompiled,
	}

	return map[string]tengo.Object{
		"write_file": &interop.AdvFunction{
			Name:    "write_file",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrArg("path"), interop.StrArg("data")},
			Value:   m.writeFile,
		},
		"write_file_lines": &interop.AdvFunction{
			Name:    "write_file_lines",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrArg("path"), interop.StrSliceArg("lines", false)},
			Value:   m.writeFileLines,
		},
		"read_file_lines": &interop.AdvFunction{
			Name:    "read_file_lines",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("path")},
			Value:   m.readFileLines,
		},
		"regex_replace_file": &interop.AdvFunction{
			Name:    "regex_replace_file",
			NumArgs: interop.ExactArgs(3),
			Args:    []interop.AdvArg{interop.StrArg("path"), interop.RegexArg("regex"), interop.StrArg("replace")},
			Value:   m.regexReplaceFile,
		},
		"mkdir_all": &interop.AdvFunction{
			Name:    "mkdir_all",
			NumArgs: interop.MinArgs(1),
			Args:    []interop.AdvArg{interop.StrSliceArg("paths", true)},
			Value:   m.mkdirAll,
		},
		"mkdir_temp": &interop.AdvFunction{
			Name:    "mkdir_temp",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrArg("dir"), interop.StrArg("pattern")},
			Value:   m.mkdirTemp,
		},
		"read_stdin": &interop.AdvFunction{
			Name:  "read_stdin",
			Value: m.readStdin,
		},
		"temp_chdir": &interop.AdvFunction{
			Name:    "temp_chdir",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrArg("path"), interop.CompileFuncArg("fn")},
			Value:   m.tempChdir,
		},
		"copy_files": &interop.AdvFunction{
			Name:    "copy_files",
			NumArgs: interop.ExactArgs(2),
			Args: []interop.AdvArg{
				interop.UnionArg("src", interop.StrType, interop.StrSliceType),
				interop.StrArg("dest"),
			},
			Value: m.copyFiles,
		},
		"copy_dirs": &interop.AdvFunction{
			Name:    "copy_dirs",
			NumArgs: interop.MinArgs(2),
			Args: []interop.AdvArg{
				interop.UnionArg("src", interop.StrType, interop.StrSliceType),
				interop.StrArg("dest"),
			},
			Value: m.copyDirs,
		},
		"prompt": &interop.AdvFunction{
			Name:    "prompt",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("msg")},
			Value:   m.promptUser,
		},
	}
}

// writeFile is like the tengo 'os.write_file' function except the file is written with 0644 permissions
// Represents 'os2.write_file(path string, data string) error'
func (m *module) writeFile(args map[string]interface{}) (tengo.Object, error) {
	path := args["path"].(string)
	data := args["data"].(string)

	err := fileutil.WriteString(path, data)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

// writeFileLines is like the writeFile function except each element in the slice is written on a new line
// Represents 'os2.write_file_lines(path string, lines []string) error'
func (m *module) writeFileLines(args map[string]interface{}) (tengo.Object, error) {
	path := args["path"].(string)
	lines := args["lines"].([]string)

	err := fileutil.WriteLines(path, lines)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

// readFileLines reads the file and splits the contents by each new line
// Represents 'os2.read_file_lines(path string) []string|error'
func (m *module) readFileLines(args map[string]interface{}) (tengo.Object, error) {
	path := args["path"].(string)

	lines, err := fileutil.ReadLines(path)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return interop.GoStrSliceToTArray(lines), nil
}

// regexReplaceFile reads the file, replaces the contents that match the regex and writes it back to the file.
// Represents 'os2.regex_replace_file(path string, regex string, replace string) error'
func (m *module) regexReplaceFile(args map[string]interface{}) (tengo.Object, error) {
	path := args["path"].(string)
	re := args["regex"].(*regexp.Regexp)
	replace := args["replace"].(string)

	data, err := os.ReadFile(path)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	replaced := re.ReplaceAll(data, []byte(replace))

	err = fileutil.WriteString(path, string(replaced))
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

// mkdirAll is a simple tengo function wrapper to 'os.MkdirAll' except it sets the directory permissions to 0755
// Represents 'os2.mkdir_all(paths ...string) error'
func (m *module) mkdirAll(args map[string]interface{}) (tengo.Object, error) {
	paths := args["paths"].([]string)
	for _, path := range paths {
		err := os.MkdirAll(path, fileutil.DefaultDirPerms)
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}
	}

	return nil, nil
}

// mkdirTemp is a tengo function wrapper to the os.MkdirTemp function
// Represents 'os2.mkdir_temp(dir string, pattern string) string|error'
func (m *module) mkdirTemp(args map[string]interface{}) (tengo.Object, error) {
	dir := args["dir"].(string)
	pattern := args["pattern"].(string)

	tempDir, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return &tengo.String{
		Value: tempDir,
	}, nil
}

// readStdin reads the current process's Stdin if anything was piped to it.
// Represents 'os2.read_stdin() []string'
func (m *module) readStdin(args map[string]interface{}) (tengo.Object, error) {
	if !fileutil.HasStdin() {
		return nil, nil
	}

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return interop.GoStrSliceToTArray(lines), nil
}

// tempChdir changes the current directory, executes the function, then changes the current directory back.
// Represents 'os2.temp_chdir(dir string, fn func())'
func (m *module) tempChdir(args map[string]interface{}) (tengo.Object, error) {
	if m.getCompiled == nil {
		return nil, errors.New("module not setup to run compiled functions from Go code")
	}

	compiled, _ := m.getCompiled()

	path := args["path"].(string)
	fn := args["fn"].(*tengo.CompiledFunction)

	var err error
	previousDir := ""

	if path != "" {
		previousDir, err = os.Getwd()
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}

		err = os.Chdir(path)
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}
	}

	runner := interop.NewCompiledFuncRunner(fn, compiled, context.Background())
	_, err = runner.Run()
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	if path != "" {
		err = os.Chdir(previousDir)
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}
	}

	return nil, nil
}

// copyFiles copies the specified files to the destination.
// Represents 'os2.copy_files(src string|[]string, dest string) error'
func (m *module) copyFiles(args map[string]interface{}) (tengo.Object, error) {
	files, ok := args["src"].([]string)
	if !ok {
		src := args["src"].(string)

		var err error
		files, err = doublestar.FilepathGlob(src)
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}
	}

	dest := args["dest"].(string)

	for _, file := range files {
		err := fileutil.CopyFile(file, dest)
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}
	}

	return nil, nil
}

// copyDirs copies the specified directories to the destination.
// Represents 'os2.copy_dirs(src string|[]string, dest string) error'
func (m *module) copyDirs(args map[string]interface{}) (tengo.Object, error) {
	srcDirs, ok := args["src"].([]string)
	if !ok {
		src := args["src"].(string)
		srcDirs = []string{src}
	}
	dest := args["dest"].(string)

	if len(srcDirs) > 1 && !fileutil.DirExists(dest) {
		return interop.GoErrToTErr(fmt.Errorf("%s: No such directory", dest)), nil
	}

	for _, src := range srcDirs {
		err := fileutil.CopyDir(src, dest)
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}
	}

	return nil, nil
}

// promptUser prints a message to Stdout and reads user input
// Represents 'os2.prompt(msg string) string|error'
func (m *module) promptUser(args map[string]interface{}) (tengo.Object, error) {
	msg := args["msg"].(string)

	fmt.Print(msg)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	err := scanner.Err()
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return interop.GoStrToTStr(scanner.Text()), nil
}
