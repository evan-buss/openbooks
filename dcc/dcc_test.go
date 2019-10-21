package dcc

import (
	"fmt"
	"testing"
)

// TestSearchParse makes sure that data is properly extracted from the DCC
// response string
func TestSearchParse(t *testing.T) {
	tables := []struct {
		search   string
		filename string
		ip       string
		port     string
		size     int
	}{
		{
			":SearchOok!ook@only.ook PRIVMSG evan_28 :DCC SEND SearchOok_results_for__hp_lovecraft.txt.zip 1543751478 2043 784",
			"SearchOok_results_for__hp_lovecraft.txt.zip", "92.3.199.54", "2043", 784},
		{
			":Search!Search@ihw-4q5hcb.dyn.suddenlink.net PRIVMSG evan_bot :DCC SEND SearchBot_results_for__stephen_king_the_stand.txt.zip 2907707975 4342 1116", "SearchBot_results_for__stephen_king_the_stand.txt.zip", "173.80.26.71", "4342", 1116},
	}

	dcc := new(Conn)
	for _, table := range tables {

		dcc.ParseSearch(table.search)
		fmt.Printf(dcc.filename)

		if dcc.filename != table.filename {
			t.Errorf("Search parser filename incorrect, got: %s, want: %s.", dcc.filename, table.filename)
		}
		if dcc.ip != table.ip {
			t.Errorf("Search parser ip incorrect, got: %s, want: %s.", dcc.ip, table.ip)
		}
		if dcc.port != table.port {
			t.Errorf("Search parser port incorrect, got: %s, want: %s.", dcc.port, table.port)
		}
		if dcc.size != table.size {
			t.Errorf("Search parser size incorrect, got: %d, want: %d.", dcc.size, table.size)
		}
	}
}
