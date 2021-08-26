package core

import (
	"reflect"
	"testing"
	"time"
)

func TestCaseInsensitiveSort(t *testing.T) {
	servers := "+FWServer ~Oatmeal +LawdyServer +fwServer"
	want := []string{"FWServer", "fwServer", "LawdyServer", "Oatmeal"}
	cache := &ServerCache{Servers: make([]string, 0), Time: time.Now()}
	cache.ParseServers(servers)

	if !reflect.DeepEqual(cache.Servers, want) {
		t.Errorf("got %#v, want %#v", cache.Servers, want)
	}
}
