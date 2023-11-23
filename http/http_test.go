package http_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/analog-substance/tengo/v2/require"
	tengohttp "github.com/analog-substance/tengomod/http"
	"github.com/analog-substance/tengomod/internal/test"
)

func expectResp(t *testing.T, expectedStatus int, expectedBody string, headers http.Header, obj interface{}) {
	require.IsType(t, &tengohttp.HTTPResponse{}, obj)

	resp := obj.(*tengohttp.HTTPResponse).Value
	require.Equal(t, resp.StatusCode, expectedStatus)

	for key, values := range headers {
		require.Equal(t, values[0], resp.Header.Get(key))
	}

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, expectedBody, string(body))
}

func testHeaders(method string) http.Header {
	return http.Header{
		"X-Header": []string{method},
	}
}

func newServer() http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Header", r.Method)

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Add("Content-Type", r.Header.Get("Content-Type"))

			body, _ := io.ReadAll(r.Body)
			fmt.Fprint(w, string(body))
		} else {
			fmt.Fprint(w, r.Method)
		}
	})

	return http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
}

func TestHTTP(t *testing.T) {
	server := newServer()
	go server.ListenAndServe()
	defer server.Close()

	u := "http://localhost:8000"
	obj := test.Module(t, "http").Call("head", u).Obj

	headers := testHeaders("HEAD")
	expectResp(t, http.StatusOK, "", headers, obj)

	obj = test.Module(t, "http").Call("get", u).Obj

	headers = testHeaders("GET")
	expectResp(t, http.StatusOK, "GET", headers, obj)

	obj = test.Module(t, "http").Call("post", u, "application/tengo", "post body").Obj

	headers = testHeaders("POST")
	headers.Add("Content-Type", "application/tengo")
	expectResp(t, http.StatusOK, "post body", headers, obj)

	obj = test.Module(t, "http").Call("put", u, "application/tengo", "put body").Obj

	headers = testHeaders("PUT")
	headers.Add("Content-Type", "application/tengo")
	expectResp(t, http.StatusOK, "put body", headers, obj)

	obj = test.Module(t, "http").Call("patch", u, "application/tengo", "patch body").Obj

	headers = testHeaders("PATCH")
	headers.Add("Content-Type", "application/tengo")
	expectResp(t, http.StatusOK, "patch body", headers, obj)

	obj = test.Module(t, "http").Call("delete", u, "application/tengo", "delete body").Obj

	headers = testHeaders("DELETE")
	headers.Add("Content-Type", "application/tengo")
	expectResp(t, http.StatusOK, "delete body", headers, obj)
}

