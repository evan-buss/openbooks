package mock

import (
	"github.com/evan-buss/openbooks/dcc"
	"log"
	"net"
	"os"
	"path/filepath"
)

type DccManager struct{}

func NewDccManager() *DccManager {
	return &DccManager{}
}

func (d *DccManager) ServeFile(path string) string {
	file, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	server := NewDccServer(path)
	selectedPortChan := make(chan int)
	go server.Start(selectedPortChan)
	log.Println(file.Size())

	dccString := dcc.New(filepath.Base(path), net.ParseIP("127.0.0.1").To4(), <-selectedPortChan, file.Size()).String()
	println(dccString)
	return dccString
}
