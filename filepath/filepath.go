package filepath

import (
	"path/filepath"
	"regexp"

	"github.com/analog-substance/fileutil"
	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengo/v2/stdlib"
	"github.com/analog-substance/tengomod/interop"
	"github.com/bmatcuk/doublestar/v4"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"join":        &tengo.UserFunction{Name: "join", Value: interop.FuncASvRS(filepath.Join)},
		"file_exists": &tengo.UserFunction{Name: "file_exists", Value: interop.FuncASRB(fileutil.FileExists)},
		"dir_exists":  &tengo.UserFunction{Name: "dir_exists", Value: interop.FuncASRB(fileutil.DirExists)},
		"base":        &tengo.UserFunction{Name: "base", Value: stdlib.FuncASRS(filepath.Base)},
		"dir":         &tengo.UserFunction{Name: "dir", Value: stdlib.FuncASRS(filepath.Dir)},
		"abs":         &tengo.UserFunction{Name: "abs", Value: stdlib.FuncASRSE(filepath.Abs)},
		"ext":         &tengo.UserFunction{Name: "ext", Value: stdlib.FuncASRS(filepath.Ext)},
		"glob":        &tengo.UserFunction{Name: "glob", Value: interop.NewCallable(glob, interop.WithArgRange(1, 2))},
		"from_slash":  &tengo.UserFunction{Name: "from_slash", Value: stdlib.FuncASRS(filepath.FromSlash)},
	}
}

func glob(args ...tengo.Object) (tengo.Object, error) {
	pattern, err := interop.TStringToGoString(args[0], "pattern")
	if err != nil {
		return nil, err
	}

	var excludeRe *regexp.Regexp
	if len(args) == 2 {
		excludePatternArg, err := interop.TStringToGoString(args[1], "exclude-pattern")
		if err != nil {
			return nil, err
		}

		excludeRe, err = regexp.Compile(excludePatternArg)
		if err != nil {
			return nil, err
		}
	}

	matches, err := doublestar.FilepathGlob(pattern)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	if excludeRe != nil {
		var filtered []string
		for _, match := range matches {
			if !excludeRe.MatchString(match) {
				filtered = append(filtered, match)
			}
		}
		return interop.GoStringSliceToTArray(filtered), nil
	}

	return interop.GoStringSliceToTArray(matches), nil
}
