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
		"write_file": &tengo.UserFunction{
			Name:  "write_file",
			Value: interop.NewCallable(m.writeFile, interop.WithExactArgs(2)),
		},
		"write_file_lines": &tengo.UserFunction{
			Name:  "write_file_lines",
			Value: interop.NewCallable(m.writeFileLines, interop.WithExactArgs(2)),
		},
		"read_file_lines": &tengo.UserFunction{
			Name:  "read_file_lines",
			Value: interop.NewCallable(m.readFileLines, interop.WithExactArgs(1)),
		},
		"regex_replace_file": &tengo.UserFunction{
			Name:  "regex_replace_file",
			Value: interop.NewCallable(m.regexReplaceFile, interop.WithExactArgs(3)),
		},
		"mkdir_all": &tengo.UserFunction{
			Name:  "mkdir_all",
			Value: interop.NewCallable(m.mkdirAll, interop.WithMinArgs(1)),
		},
		"mkdir_temp": &tengo.UserFunction{
			Name:  "mkdir_temp",
			Value: interop.NewCallable(m.mkdirTemp, interop.WithExactArgs(2)),
		},
		"read_stdin": &tengo.UserFunction{
			Name:  "read_stdin",
			Value: interop.NewCallable(m.readStdin),
		},
		"temp_chdir": &tengo.UserFunction{
			Name:  "temp_chdir",
			Value: interop.NewCallable(m.tempChdir, interop.WithExactArgs(2)),
		},
		"copy_files": &tengo.UserFunction{
			Name:  "copy_files",
			Value: interop.NewCallable(m.copyFiles, interop.WithExactArgs(2)),
		},
		"copy_dirs": &tengo.UserFunction{
			Name:  "copy_dirs",
			Value: interop.NewCallable(m.copyDirs, interop.WithMinArgs(2)),
		},
		"prompt": &tengo.UserFunction{
			Name:  "prompt",
			Value: interop.NewCallable(m.promptUser, interop.WithExactArgs(1)),
		},
	}
}

// writeFile is like the tengo 'os.write_file' function except the file is written with 0644 permissions
// Represents 'os2.write_file(path string, data string) error'
func (m *module) writeFile(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	path, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "path",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	data, ok := tengo.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "data",
			Expected: "string",
			Found:    args[1].TypeName(),
		}
	}

	err := fileutil.WriteString(path, data)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

// writeFileLines is like the writeFile function except each element in the slice is written on a new line
// Represents 'os2.write_file_lines(path string, lines []string) error'
func (m *module) writeFileLines(args ...tengo.Object) (tengo.Object, error) {
	path, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "path",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	lines, err := interop.TArrayToGoStrSlice(args[1], "lines")
	if err != nil {
		return nil, err
	}

	err = fileutil.WriteLines(path, lines)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

// regexReplaceFile reads the file, replaces the contents that match the regex and writes it back to the file.
// Represents 'os2.regex_replace_file(path string, regex string, replace string) error'
func (m *module) regexReplaceFile(args ...tengo.Object) (tengo.Object, error) {
	path, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "path",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	regex, ok := tengo.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "regex",
			Expected: "string",
			Found:    args[1].TypeName(),
		}
	}

	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	replace, ok := tengo.ToString(args[2])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "replace",
			Expected: "string",
			Found:    args[2].TypeName(),
		}
	}

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
func (m *module) mkdirAll(args ...tengo.Object) (tengo.Object, error) {
	for _, obj := range args {
		path, _ := tengo.ToString(obj)
		err := os.MkdirAll(path, fileutil.DefaultDirPerms)
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}
	}

	return nil, nil
}

// mkdirTemp is a tengo function wrapper to the os.MkdirTemp function
// Represents 'os2.mkdir_temp(dir string, pattern string) string|error'
func (m *module) mkdirTemp(args ...tengo.Object) (tengo.Object, error) {
	dir, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "dir",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	pattern, ok := tengo.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "pattern",
			Expected: "string",
			Found:    args[1].TypeName(),
		}
	}

	tempDir, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return &tengo.String{
		Value: tempDir,
	}, nil
}

// readFileLines reads the file and splits the contents by each new line
// Represents 'os2.read_file_lines(path string) []string|error'
func (m *module) readFileLines(args ...tengo.Object) (tengo.Object, error) {
	path, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "path",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	lines, err := fileutil.ReadLines(path)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return interop.GoStrSliceToTArray(lines), nil
}

// readStdin reads the current process's Stdin if anything was piped to it.
// Represents 'os2.read_stdin() []string'
func (m *module) readStdin(args ...tengo.Object) (tengo.Object, error) {
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
func (m *module) tempChdir(args ...tengo.Object) (tengo.Object, error) {
	if m.getCompiled == nil {
		return nil, errors.New("module not setup to run compiled functions from Go code")
	}

	compiled, _ := m.getCompiled()

	path, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "path",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}

	fn, ok := args[1].(*tengo.CompiledFunction)
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "fn",
			Expected: "function",
			Found:    args[1].TypeName(),
		}
	}

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
func (m *module) copyFiles(args ...tengo.Object) (tengo.Object, error) {
	files, err := interop.TArrayToGoStrSlice(args[0], "src")
	if err == nil {
		src, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "src",
				Expected: "string|[]string",
				Found:    args[0].TypeName(),
			}
		}

		files, err = doublestar.FilepathGlob(src)
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}
	}

	dest, ok := tengo.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "dest",
			Expected: "string",
			Found:    args[1].TypeName(),
		}
	}

	for _, file := range files {
		err = fileutil.CopyFile(file, dest)
		if err != nil {
			return interop.GoErrToTErr(err), nil
		}
	}

	return nil, nil
}

// copyDirs copies the specified directories to the destination.
// Represents 'os2.copy_dirs(src string..., dest string) error'
func (m *module) copyDirs(args ...tengo.Object) (tengo.Object, error) {
	var srcDirs []string
	for _, arg := range args[:len(args)-1] {
		src, ok := tengo.ToString(arg)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "src",
				Expected: "string",
				Found:    arg.TypeName(),
			}
		}

		srcDirs = append(srcDirs, src)
	}

	dest, ok := tengo.ToString(args[len(args)-1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "dest",
			Expected: "string",
			Found:    args[len(args)-1].TypeName(),
		}
	}

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
func (m *module) promptUser(args ...tengo.Object) (tengo.Object, error) {
	msg, err := interop.TStrToGoStr(args[0], "msg")
	if err != nil {
		return nil, err
	}

	fmt.Print(msg)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	err = scanner.Err()
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return interop.GoStrToTStr(scanner.Text()), nil
}
