package dcc

import "testing"

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
			"SearchOok_results_for__hp_lovecraft.txt.zip", "1543751478", "2043", 784},
		{
			":Search!Search@ihw-4q5hcb.dyn.suddenlink.net PRIVMSG evan_bot :DCC SEND SearchBot_results_for__stephen_king_the_stand.txt.zip 2907707975 4342 1116", "SearchBot_results_for__stephen_king_the_stand.txt.zip", "2907707975", "4342", 1116},
	}

	for _, table := range tables {
		filename, ip, port, size := parseSearch(table.search)
		if filename != table.filename {
			t.Errorf("Search parser filename incorrect, got: %s, want: %s.", filename, table.filename)
		}
		if ip != table.ip {
			t.Errorf("Search parser ip incorrect, got: %s, want: %s.", ip, table.ip)
		}
		if port != table.port {
			t.Errorf("Search parser port incorrect, got: %s, want: %s.", port, table.port)
		}
		if size != table.size {
			t.Errorf("Search parser size incorrect, got: %d, want: %d.", size, table.size)
		}
	}
}
