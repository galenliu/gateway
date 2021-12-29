package main

import (
	"fmt"
	"github.com/galenliu/gateway/cmd/gateway/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
