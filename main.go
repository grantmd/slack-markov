package main

// Main entry point for the app. Handles command-line options, starts the web
// listener and any import, etc

import (
	"flag"
	"fmt"
	"os"
)

var (
	httpPort  int
	numWords  int
	prefixLen int

	markovChain *Chain
)

func main() {
	// Parse command-line options
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: ./slack-markov -port=8000\n")
		flag.PrintDefaults()
	}

	flag.IntVar(&httpPort, "port", 8000, "The HTTP port on which to listen")
	flag.IntVar(&numWords, "words", 100, "Maximum number of words in the output")
	flag.IntVar(&prefixLen, "prefix", 2, "Prefix length in words")

	flag.Parse()

	if httpPort == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Rebuild the markov chain from state
	markovChain = NewChain(prefixLen) // Initialize a new Chain.

	// Start the webserver
	StartServer(httpPort)
}
