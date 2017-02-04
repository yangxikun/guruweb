package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
    "github.com/yangxikun/guruweb/output"
)

var (
	httpAddr = flag.String("http", ":8080", "HTTP listen address")
	verbose  = flag.Bool("v", false, "verbose mode")
	open     = flag.Bool("open", true, "try to open browser")
	scope    = flag.String("scope", "...", "comma-separated list of `packages` the analysis should be limited to.\n\t" +
                    "see buildutil.ExpandPatterns for more help")
)

func main() {

	flag.Parse()
	initConfig(strings.Split(*scope, ","), *verbose)
	err := initIndex()
	exitOn(err)

	registerHandlers()

	srv := &http.Server{Addr: *httpAddr}
	l, err := net.Listen("tcp", srv.Addr)
	exitOn(err)
	url := fmt.Sprintf("http://localhost%s/", *httpAddr)
	if *open {
		err = startBrowser(url)
		if err != nil {
			output.Warn(err.Error())
			output.Warn("open browser fail")
		}
	}
	output.Normal("GOPATH: %s", os.Getenv("GOPATH"))
	output.Normal("GuruWeb is running at %s", url)
	exitError(srv.Serve(l))
}

// startBrowser tries to open the URL in a browser
// and reports whether it succeeds.
func startBrowser(url string) error {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start()
}

func registerHandlers() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/recommend-search", serveRecommendSearch)
	http.HandleFunc("/recommend-pkgs", serveRecommendPkgs)
	http.HandleFunc("/file", serveFile)
	http.HandleFunc("/query", serveQuery)
	http.HandleFunc("/config", serveConfig)
	staticPrefix := "/static/"
	http.Handle(staticPrefix, http.StripPrefix(staticPrefix, http.HandlerFunc(ServeStatic)))
}

func exitOn(err error) {
	if err != nil {
		exitError(err)
	}
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
