package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/DhanushNehru/hypersweep/pkg/checker"
	"github.com/DhanushNehru/hypersweep/pkg/extractor"
	"github.com/DhanushNehru/hypersweep/pkg/reporter"
)

func main() {
	pathPtr := flag.String("path", ".", "Path to scan for URLs (defaults to current directory)")
	workersPtr := flag.Int("workers", 50, "Number of concurrent workers")
	timeoutPtr := flag.Int("timeout", 10, "Timeout for each HTTP request in seconds")
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "HyperSweep - Blazing fast concurrent link checker\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	fmt.Printf("Starting HyperSweep on path: %s (Workers: %d)\n", *pathPtr, *workersPtr)
	startTime := time.Now()

	// 1. Extract URLs
	ext := extractor.NewExtractor(*pathPtr)
	urls, err := ext.Extract()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting URLs: %v\n", err)
		os.Exit(1)
	}

	if len(urls) == 0 {
		fmt.Println("No URLs found.")
		os.Exit(0)
	}
	fmt.Printf("Extracted %d URLs. Verifying...\n", len(urls))

	// 2. Check URLs Concurrently
	chk_timeout := time.Duration(*timeoutPtr) * time.Second
	chk := checker.NewChecker(*workersPtr, chk_timeout)
	results := chk.Run(urls)

	// 3. Report Results
	hasBroken := reporter.PrintResults(results, time.Since(startTime))

	// Exit with code 1 if broken links found (Crucial for CI/CD)
	if hasBroken {
		os.Exit(1)
	}
}
