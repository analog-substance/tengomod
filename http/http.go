package http

import (
	"net/http"
	"net/url"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

var defaultClient *HTTPClient = makeHTTPClient(http.DefaultClient)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"method_get": &tengo.String{
			Value: http.MethodGet,
		},
		"method_put": &tengo.String{
			Value: http.MethodPut,
		},
		"method_post": &tengo.String{
			Value: http.MethodPost,
		},
		"method_delete": &tengo.String{
			Value: http.MethodDelete,
		},
		"method_head": &tengo.String{
			Value: http.MethodHead,
		},
		"method_options": &tengo.String{
			Value: http.MethodOptions,
		},
		"method_patch": &tengo.String{
			Value: http.MethodPatch,
		},
		"default_client": defaultClient,
		"head": &interop.AdvFunction{
			Name:    "head",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("url")},
			Value:   defaultClient.head,
		},
		"get": &interop.AdvFunction{
			Name:    "get",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("url")},
			Value:   defaultClient.get,
		},
		"post": &interop.AdvFunction{
			Name:    "post",
			NumArgs: interop.ArgRange(1, 3),
			Args: []interop.AdvArg{
				interop.StrArg("url"),
				interop.StrArg("contentType"),
				interop.UnionArg("body", interop.ByteSliceType, interop.ObjectType),
			},
			Value: defaultClient.post,
		},
		"put": &interop.AdvFunction{
			Name:    "put",
			NumArgs: interop.ArgRange(1, 3),
			Args: []interop.AdvArg{
				interop.StrArg("url"),
				interop.StrArg("contentType"),
				interop.UnionArg("body", interop.ByteSliceType, interop.ObjectType),
			},
			Value: defaultClient.put,
		},
		"patch": &interop.AdvFunction{
			Name:    "patch",
			NumArgs: interop.ArgRange(1, 3),
			Args: []interop.AdvArg{
				interop.StrArg("url"),
				interop.StrArg("contentType"),
				interop.UnionArg("body", interop.ByteSliceType, interop.ObjectType),
			},
			Value: defaultClient.patch,
		},
		"delete": &interop.AdvFunction{
			Name:    "delete",
			NumArgs: interop.ArgRange(1, 3),
			Args: []interop.AdvArg{
				interop.StrArg("url"),
				interop.StrArg("contentType"),
				interop.UnionArg("body", interop.ByteSliceType, interop.ObjectType),
			},
			Value: defaultClient.delete,
		},
		"new_client": &interop.AdvFunction{
			Name:    "new_client",
			NumArgs: interop.MaxArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("baseURL")},
			Value:   newHTTPClient,
		},
		"new_request": &interop.AdvFunction{
			Name:    "new_request",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrArg("method"), interop.StrArg("url")},
			Value:   newRequest,
		},
	}
}

func newHTTPClient(args interop.ArgMap) (tengo.Object, error) {
	baseURL, _ := args.GetString("baseURL")

	client := makeHTTPClient(&http.Client{})
	client.SetBaseURL(baseURL)

	return client, nil
}

func newRequest(args interop.ArgMap) (tengo.Object, error) {
	method, _ := args.GetString("method")
	u, _ := args.GetString("url")

	reqURL, err := url.Parse(u)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	req := &http.Request{
		Method: method,
		URL:    reqURL,
		Header: make(http.Header),
	}

	return makeHTTPRequest(req), nil
}
