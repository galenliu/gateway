package main

import (
	"flag"
	"fmt"
	"github.com/galenliu/gateway/cmd/gateway/cmd"
	"os"
	"os/signal"
	"syscall"
)

var (
	proFile     string
	showVersion bool
)

func init() {
	flag.StringVar(&proFile, "profile", "", "Profile directory")
	flag.BoolVar(&showVersion, "version", false, "version")
}



// TermFunc defines the function which is executed on termination.
type TermFunc func(sig os.Signal)

// OnTermination calls a function when the app receives an interrupt of kill signal.
func OnTermination(fn TermFunc) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGTERM)

	func() {
		select {
		case sig := <-c:
			if fn != nil {
				fn(sig)
			}
		}
	}()
}

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
