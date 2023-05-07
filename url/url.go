package url

import (
	"net/url"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"hostname": &tengo.UserFunction{
			Name:  "hostname",
			Value: interop.NewCallable(hostname, interop.WithExactArgs(1)),
		},
	}
}

func hostname(args ...tengo.Object) (tengo.Object, error) {
	rawURL, err := interop.TStrToGoStr(args[0], "url")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return interop.GoStrToTStr(parsedURL.Hostname()), nil
}
