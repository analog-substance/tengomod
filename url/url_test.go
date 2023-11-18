package url

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

func Test_hostname(t *testing.T) {
	withHostname, _ := url.Parse("http://example.com")
	noHostname, _ := url.Parse("example.com")

	type args struct {
		args interop.ArgMap
	}
	tests := []struct {
		name string
		args args
		want tengo.Object
	}{
		{
			name: "Full URL with Hostname",
			args: args{
				args: interop.ArgMap{
					"url": withHostname,
				},
			},
			want: &tengo.String{Value: "example.com"},
		},
		{
			name: "Not a Full URL",
			args: args{
				args: interop.ArgMap{
					"url": noHostname,
				},
			},
			want: &tengo.String{Value: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := hostname(tt.args.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hostname() = %v, want %v", got, tt.want)
			}
		})
	}
}
