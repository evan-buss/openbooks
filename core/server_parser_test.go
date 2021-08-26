package core

import (
	"reflect"
	"testing"
	"time"
)

func TestCaseInsensitiveSort(t *testing.T) {
	cache := &ServerCache{Servers: make([]string, 0), Time: time.Now()}
	cases := []struct {
		input string
		want  []string
	}{
		{"+FWServer ~Oatmeal +LawdyServer +fwServer", []string{"FWServer", "fwServer", "LawdyServer", "Oatmeal"}},
		{"", []string{}},
	}

	for _, v := range cases {
		cache.ParseServers(v.input)
		if !reflect.DeepEqual(cache.Servers, v.want) {
			t.Errorf("got %#v, want %#v", cache.Servers, v.want)
		}
	}
}
