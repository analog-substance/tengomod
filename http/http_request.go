package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengo/v2/stdlib/json"
	"github.com/analog-substance/tengomod/interop"
	"github.com/analog-substance/tengomod/types"
)

type HTTPRequest struct {
	types.PropObject
	Value *http.Request
}

func (r *HTTPRequest) TypeName() string {
	return "http-request"
}

// String should return a string representation of the type's value.
func (r *HTTPRequest) String() string {
	return "<http-request>"
}

// IsFalsy should return true if the value of the type should be considered
// as falsy.
func (r *HTTPRequest) IsFalsy() bool {
	return r.Value == nil
}

// CanIterate should return whether the Object can be Iterated.
func (r *HTTPRequest) CanIterate() bool {
	return false
}

func (r *HTTPRequest) token(args interop.ArgMap) (tengo.Object, error) {
	token, _ := args.GetString("token")

	r.Value.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil, nil
}

func (r *HTTPRequest) userAgent(args interop.ArgMap) (tengo.Object, error) {
	userAgent, _ := args.GetString("userAgent")

	r.Value.Header.Set("User-Agent", userAgent)
	return nil, nil
}

func (r *HTTPRequest) contentType(args interop.ArgMap) (tengo.Object, error) {
	contentType, _ := args.GetString("contentType")

	r.Value.Header.Set("Content-Type", contentType)
	return nil, nil
}

func (r *HTTPRequest) clone(_ ...tengo.Object) (tengo.Object, error) {
	req := r.Value.Clone(context.Background())

	// Make 2 copies of body
	body, _ := io.ReadAll(r.Value.Body)
	r.Value.Body = io.NopCloser(bytes.NewBuffer(body))
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	return makeHTTPRequest(req), nil
}

func makeHTTPRequest(r *http.Request) *HTTPRequest {
	request := &HTTPRequest{
		Value: r,
	}

	objectMap := map[string]tengo.Object{
		"user_agent": &interop.AdvFunction{
			Name:    "user_agent",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("userAgent")},
			Value:   request.userAgent,
		},
		"token": &interop.AdvFunction{
			Name:    "token",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("token")},
			Value:   request.token,
		},
		"content_type": &interop.AdvFunction{
			Name:    "content_type",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("contentType")},
			Value:   request.contentType,
		},
		"clone": &tengo.UserFunction{
			Name:  "clone",
			Value: request.clone,
		},
	}
	properties := map[string]types.Property{
		"method": {
			Get: func() tengo.Object {
				return interop.GoStrToTStr(request.Value.Method)
			},
			Set: func(o tengo.Object) error {
				method, err := interop.TStrToGoStr(o, "method")
				if err != nil {
					return err
				}

				request.Value.Method = method
				return nil
			},
		},
		"header": {
			Get: func() tengo.Object {
				return makeHTTPHeader(request.Value.Header)
			},
		},
		"url": {
			Get: func() tengo.Object {
				return interop.GoStrToTStr(request.Value.URL.String())
			},
			Set: func(o tengo.Object) error {
				rawURL, err := interop.TStrToGoStr(o, "url")
				if err != nil {
					return err
				}

				u, err := url.Parse(rawURL)
				if err != nil {
					return err
				}

				request.Value.URL = u
				return nil
			},
		},
		"body": {
			Get: func() tengo.Object {
				body, _ := io.ReadAll(request.Value.Body)
				return &tengo.Bytes{
					Value: body,
				}
			},
			Set: func(o tengo.Object) error {
				body, ok := tengo.ToByteSlice(o)
				if !ok {
					var err error
					body, err = json.Encode(o)
					if err != nil {
						return err
					}
				}

				request.Value.ContentLength = int64(len(body))
				request.Value.Body = io.NopCloser(bytes.NewBuffer(body))
				return nil
			},
		},
	}

	request.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: properties,
	}

	return request
}
