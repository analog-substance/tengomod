package url

import (
	"net/url"

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
func hostname(args map[string]interface{}) (tengo.Object, error) {
	parsedURL := args["url"].(*url.URL)
	return interop.GoStrToTStr(parsedURL.Hostname()), nil
}
