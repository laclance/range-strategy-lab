package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"range-strategy-lab/internal/lab"
)

const edgeFadeFlagName = "backtest-first-btc-15m-range-edge-exhaustion-fade-v1"

func init() {
	if !edgeFadeFlagPresent(os.Args[1:]) {
		return
	}
	if err := runEdgeFadeWithArgs(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func edgeFadeFlagPresent(args []string) bool {
	prefix := "-" + edgeFadeFlagName
	for _, arg := range args {
		if arg == prefix || strings.HasPrefix(arg, prefix+"=") {
			return true
		}
	}
	return false
}
