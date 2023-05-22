package url

import (
	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"hostname": &interop.AdvFunction{
			Name:    "hostname",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.URLArg("url")},
			Value:   hostname,
		},
	}
}

// hostname returns the hostname of the URL
// Represents 'url.hostname(url string) string|error'
func hostname(args interop.ArgMap) (tengo.Object, error) {
	parsedURL, _ := args.GetURL("url")
	return interop.GoStrToTStr(parsedURL.Hostname()), nil
}
