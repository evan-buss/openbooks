package core

// Config contains user options used throughout application
type Config struct {
	CliMode     bool
	Log         bool
	OpenBrowser bool
	Port        string
	UserName    string
	Persist     bool
	DownloadDir string
}
