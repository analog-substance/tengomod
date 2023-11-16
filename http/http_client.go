package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengo/v2/stdlib/json"
	"github.com/analog-substance/tengomod/interop"
	"github.com/analog-substance/tengomod/types"
)

type HTTPClient struct {
	types.PropObject
	Value   *http.Client
	baseURL string
	header  http.Header
}

func (c *HTTPClient) TypeName() string {
	return "http-client"
}

// String should return a string representation of the type's value.
func (c *HTTPClient) String() string {
	return "<http-client>"
}

// IsFalsy should return true if the value of the type should be considered
// as falsy.
func (c *HTTPClient) IsFalsy() bool {
	return c.Value == nil
}

func (c *HTTPClient) SetBaseURL(u string) {
	c.baseURL = strings.TrimRight(u, "/")
}

func (c *HTTPClient) transport() *http.Transport {
	if c.Value.Transport == nil {
		c.Value.Transport = &http.Transport{}
	}

	return c.Value.Transport.(*http.Transport)
}

func (c *HTTPClient) tlsConfig() *tls.Config {
	transport := c.transport()
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}

	return transport.TLSClientConfig
}

func (c *HTTPClient) newRequest(method string, u string) (*http.Request, error) {
	var err error
	var reqURL *url.URL
	if c.baseURL != "" {
		reqURL, err = url.Parse(fmt.Sprintf("%s/%s", c.baseURL, strings.TrimLeft(u, "/")))
	} else {
		reqURL, err = url.Parse(u)
	}

	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: method,
		URL:    reqURL,
		Header: c.header.Clone(),
	}

	return req, nil
}

func (c *HTTPClient) newBodyRequest(method string, args interop.ArgMap) (*http.Request, error) {
	u, _ := args.GetString("url")
	contentType, _ := args.GetString("contentType")

	body, err := c.getBodyArg(args)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(method, u)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
}

func (c *HTTPClient) getBodyArg(args interop.ArgMap) ([]byte, error) {
	var body []byte
	if args.Exists("body") {
		var ok bool
		body, ok = args.GetByteSlice("body")
		if !ok {
			obj, _ := args.GetObject("body")

			var err error
			body, err = json.Encode(obj)
			if err != nil {
				return nil, err
			}
		}
	}

	return body, nil
}

func (c *HTTPClient) do(req *http.Request) (tengo.Object, error) {
	resp, err := c.Value.Do(req)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return makeHTTPResponse(resp), nil
}

func (c *HTTPClient) tengoDo(args interop.ArgMap) (tengo.Object, error) {
	obj, _ := args.GetObject("request")
	req := obj.(*HTTPRequest)

	return c.do(req.Value)
}

func (c *HTTPClient) tengoNewRequest(args interop.ArgMap) (tengo.Object, error) {
	method, _ := args.GetString("method")
	u, _ := args.GetString("url")

	req, err := c.newRequest(method, u)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return makeHTTPRequest(req), nil
}

func (c *HTTPClient) head(args interop.ArgMap) (tengo.Object, error) {
	u, _ := args.GetString("url")

	req, err := c.newRequest(http.MethodHead, u)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return c.do(req)
}

func (c *HTTPClient) get(args interop.ArgMap) (tengo.Object, error) {
	u, _ := args.GetString("url")

	req, err := c.newRequest(http.MethodGet, u)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return c.do(req)
}

func (c *HTTPClient) post(args interop.ArgMap) (tengo.Object, error) {
	req, err := c.newBodyRequest(http.MethodPost, args)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return c.do(req)
}

func (c *HTTPClient) put(args interop.ArgMap) (tengo.Object, error) {
	req, err := c.newBodyRequest(http.MethodPut, args)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return c.do(req)
}

func (c *HTTPClient) patch(args interop.ArgMap) (tengo.Object, error) {
	req, err := c.newBodyRequest(http.MethodPatch, args)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return c.do(req)
}

func (c *HTTPClient) delete(args interop.ArgMap) (tengo.Object, error) {
	req, err := c.newBodyRequest(http.MethodDelete, args)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return c.do(req)
}

func (c *HTTPClient) proxy(args interop.ArgMap) (tengo.Object, error) {
	proxyURL, _ := args.GetURL("url")

	c.transport().Proxy = http.ProxyURL(proxyURL)
	return c, nil
}

func (c *HTTPClient) userAgent(args interop.ArgMap) (tengo.Object, error) {
	userAgent, _ := args.GetString("userAgent")

	c.header.Set("User-Agent", userAgent)
	return c, nil
}

func (c *HTTPClient) token(args interop.ArgMap) (tengo.Object, error) {
	token, _ := args.GetString("token")

	c.header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return c, nil
}

func (c *HTTPClient) disableRedirects(args ...tengo.Object) (tengo.Object, error) {
	c.Value.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return c, nil
}

func (c *HTTPClient) insecureSkipVerify(args ...tengo.Object) (tengo.Object, error) {
	c.tlsConfig().InsecureSkipVerify = true
	return c, nil
}

func makeHTTPClient(c *http.Client) *HTTPClient {
	client := &HTTPClient{
		Value:  c,
		header: make(http.Header),
	}

	objectMap := map[string]tengo.Object{
		"do": &interop.AdvFunction{
			Name:    "do",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.CustomArg("request", &HTTPRequest{})},
			Value:   client.tengoDo,
		},
		"head": &interop.AdvFunction{
			Name:    "head",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("url")},
			Value:   client.head,
		},
		"get": &interop.AdvFunction{
			Name:    "get",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("url")},
			Value:   client.get,
		},
		"post": &interop.AdvFunction{
			Name:    "post",
			NumArgs: interop.ArgRange(1, 3),
			Args: []interop.AdvArg{
				interop.StrArg("url"),
				interop.StrArg("contentType"),
				interop.UnionArg("body", interop.ByteSliceType, interop.ObjectType),
			},
			Value: client.post,
		},
		"put": &interop.AdvFunction{
			Name:    "put",
			NumArgs: interop.ArgRange(1, 3),
			Args: []interop.AdvArg{
				interop.StrArg("url"),
				interop.StrArg("contentType"),
				interop.UnionArg("body", interop.ByteSliceType, interop.ObjectType),
			},
			Value: client.put,
		},
		"patch": &interop.AdvFunction{
			Name:    "patch",
			NumArgs: interop.ArgRange(1, 3),
			Args: []interop.AdvArg{
				interop.StrArg("url"),
				interop.StrArg("contentType"),
				interop.UnionArg("body", interop.ByteSliceType, interop.ObjectType),
			},
			Value: client.patch,
		},
		"delete": &interop.AdvFunction{
			Name:    "delete",
			NumArgs: interop.ArgRange(1, 3),
			Args: []interop.AdvArg{
				interop.StrArg("url"),
				interop.StrArg("contentType"),
				interop.UnionArg("body", interop.ByteSliceType, interop.ObjectType),
			},
			Value: client.delete,
		},
		"new_request": &interop.AdvFunction{
			Name:    "new_request",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrArg("method"), interop.StrArg("url")},
			Value:   client.tengoNewRequest,
		},
		"proxy": &interop.AdvFunction{
			Name:    "proxy",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.URLArg("url")},
			Value:   client.proxy,
		},
		"user_agent": &interop.AdvFunction{
			Name:    "user_agent",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("userAgent")},
			Value:   client.userAgent,
		},
		"token": &interop.AdvFunction{
			Name:    "token",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("token")},
			Value:   client.token,
		},
		"disable_redirects": &tengo.UserFunction{
			Name:  "disable_redirects",
			Value: client.disableRedirects,
		},
		"insecure_skip_verify": &tengo.UserFunction{
			Name:  "insecure_skip_verify",
			Value: client.insecureSkipVerify,
		},
	}

	properties := map[string]types.Property{
		"header": {
			Get: func() tengo.Object {
				return makeHTTPHeader(client.header)
			},
		},
		"base_url": {
			Get: func() tengo.Object {
				return interop.GoStrToTStr(client.baseURL)
			},
			Set: func(o tengo.Object) error {
				u, err := interop.TStrToGoStr(o, "base_url")
				if err != nil {
					return err
				}

				client.SetBaseURL(u)
				return nil
			},
		},
	}

	client.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: properties,
	}

	return client
}
