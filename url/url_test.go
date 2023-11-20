package url_test

import (
	"testing"

	"github.com/analog-substance/tengomod/internal/test"
)

func TestURL(t *testing.T) {
	test.Module(t, "url").Call("hostname", "http://example.com").Expect("example.com")
	test.Module(t, "url").Call("hostname", "example.com").Expect("")
}
