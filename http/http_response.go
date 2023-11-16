package http

import (
	"io"
	"net/http"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengo/v2/stdlib/json"
	"github.com/analog-substance/tengomod/interop"
	"github.com/analog-substance/tengomod/types"
)

type HTTPResponse struct {
	types.PropObject
	Value *http.Response
	body  []byte
}

func (r *HTTPResponse) TypeName() string {
	return "http-response"
}

// String should return a string representation of the type's value.
func (r *HTTPResponse) String() string {
	return "<http-response>"
}

// IsFalsy should return true if the value of the type should be considered
// as falsy.
func (r *HTTPResponse) IsFalsy() bool {
	return r.Value == nil
}

func (r *HTTPResponse) ensureBody() ([]byte, error) {
	if len(r.body) == 0 {
		body, err := io.ReadAll(r.Value.Body)
		if err != nil && len(body) == 0 {
			return nil, err
		}

		r.body = body
	}

	return r.body, nil
}

func (r *HTTPResponse) getBody(args ...tengo.Object) (tengo.Object, error) {
	body, err := r.ensureBody()
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return &tengo.Bytes{
		Value: body,
	}, nil
}

func (r *HTTPResponse) unmarshalJSON(args ...tengo.Object) (tengo.Object, error) {
	body, err := r.ensureBody()
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	decoded, err := json.Decode(body)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return decoded, nil
}

func (r *HTTPResponse) isErrorCode(args ...tengo.Object) (tengo.Object, error) {
	isError := r.Value.StatusCode >= 400 && r.Value.StatusCode < 600
	return interop.GoBoolToTBool(isError), nil
}

func (r *HTTPResponse) isSuccessCode(args ...tengo.Object) (tengo.Object, error) {
	isSuccess := r.Value.StatusCode >= 200 && r.Value.StatusCode < 300
	return interop.GoBoolToTBool(isSuccess), nil
}

func (r *HTTPResponse) isRedirect(args ...tengo.Object) (tengo.Object, error) {
	isRedirect := r.Value.StatusCode >= 300 && r.Value.StatusCode < 400
	return interop.GoBoolToTBool(isRedirect), nil
}

func makeHTTPResponse(r *http.Response) *HTTPResponse {
	response := &HTTPResponse{
		Value: r,
	}

	objectMap := map[string]tengo.Object{
		"body": &tengo.UserFunction{
			Name:  "body",
			Value: response.getBody,
		},
		"unmarshal_json": &tengo.UserFunction{
			Name:  "unmarshal_json",
			Value: response.unmarshalJSON,
		},
		"is_error": &tengo.UserFunction{
			Name:  "is_error",
			Value: response.isErrorCode,
		},
		"is_success": &tengo.UserFunction{
			Name:  "is_success",
			Value: response.isSuccessCode,
		},
		"is_redirect": &tengo.UserFunction{
			Name:  "is_redirect",
			Value: response.isRedirect,
		},
	}

	properties := map[string]types.Property{
		"header": {
			Get: func() tengo.Object {
				return makeHTTPHeader(response.Value.Header)
			},
		},
		"content_length": types.StaticProperty(interop.GoIntToTInt(int(r.ContentLength))),
		"status":         types.StaticProperty(interop.GoStrToTStr(r.Status)),
		"status_code":    types.StaticProperty(interop.GoIntToTInt(r.StatusCode)),
		"request":        types.StaticProperty(makeHTTPRequest(r.Request)),
	}

	response.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: properties,
	}

	return response
}
