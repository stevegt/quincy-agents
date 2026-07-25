package main

import (
	"fmt"
	"os"

	"github.com/Qu1ncyRy4n/Agents/tools/mogent/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
