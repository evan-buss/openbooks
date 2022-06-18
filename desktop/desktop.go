//go:build !webview

package desktop

import "github.com/evan-buss/openbooks/util"

func StartWebView(url string, debug bool) {
	util.OpenBrowser(url)
}
