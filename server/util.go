package server

import (
	"log"
	"os/exec"
	"runtime"
)

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		log.Println("Error opening browser. Unsupported platform.")
	}
	if err != nil {
		log.Println(err)
	}
}
