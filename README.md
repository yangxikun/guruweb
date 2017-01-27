guruweb is a browser based user interface for the Go tool guru, which is a source code comprehension tool for Go programs.
Inspired by [fzipp/pythia](https://github.com/fzipp/pythia), which use Go tool oracle(replaced by guru).

For more information on the Go guru, see the [original announcement](https://godoc.org/golang.org/x/tools/cmd/guru).

Installing from source
----------------------

Building Pythia requires at least Go version 1.2 or higher.

To install, run

    $ go get github.com/yangxikun/guruweb

You will now find a `guruweb` binary in your `$GOPATH/bin` directory.

Usage
-----

Start the web application with a package path, e.g.:

    $ guruweb -scope net/http

By default it will listen on port :8080 and try to launch the application
in your browser. You can choose a different port via the `-http` flag, e.g.:

    $ guruweb -http :6060 -scope fmt

Run `guruweb -help` for more information.
