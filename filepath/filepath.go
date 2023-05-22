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
		"join": &tengo.UserFunction{
			Name:  "join",
			Value: interop.FuncASvRS(filepath.Join),
		},
		"file_exists": &tengo.UserFunction{
			Name:  "file_exists",
			Value: interop.FuncASRB(fileutil.FileExists),
		},
		"dir_exists": &tengo.UserFunction{
			Name:  "dir_exists",
			Value: interop.FuncASRB(fileutil.DirExists),
		},
		"base": &tengo.UserFunction{
			Name:  "base",
			Value: stdlib.FuncASRS(filepath.Base),
		},
		"dir": &tengo.UserFunction{
			Name:  "dir",
			Value: stdlib.FuncASRS(filepath.Dir),
		},
		"abs": &tengo.UserFunction{
			Name:  "abs",
			Value: stdlib.FuncASRSE(filepath.Abs),
		},
		"ext": &tengo.UserFunction{
			Name:  "ext",
			Value: stdlib.FuncASRS(filepath.Ext),
		},
		"glob": &interop.AdvFunction{
			Name:    "glob",
			NumArgs: interop.ArgRange(1, 2),
			Args:    []interop.AdvArg{interop.StrArg("pattern"), interop.RegexArg("exclude-pattern")},
			Value:   glob,
		},
		"from_slash": &tengo.UserFunction{
			Name:  "from_slash",
			Value: stdlib.FuncASRS(filepath.FromSlash),
		},
	}
}

func glob(args map[string]interface{}) (tengo.Object, error) {
	pattern := args["pattern"].(string)

	var excludeRe *regexp.Regexp
	if value, ok := args["exclude-pattern"]; ok {
		excludeRe = value.(*regexp.Regexp)
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
		return interop.GoStrSliceToTArray(filtered), nil
	}

	return interop.GoStrSliceToTArray(matches), nil
}
