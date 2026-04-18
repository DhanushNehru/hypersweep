package reporter

import (
	"fmt"
	"time"

	"github.com/DhanushNehru/hypersweep/pkg/checker"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// PrintResults formats and prints the results to stdout
// It returns true if there were any broken links
func PrintResults(results []checker.CheckResult, duration time.Duration) bool {
	var broken []checker.CheckResult
	var successCount int

	fmt.Printf("\n%s--- HyperSweep Results ---%s\n", colorBold, colorReset)

	for _, res := range results {
		if res.IsAlive {
			successCount++
			// Optional: print successful links (can be noisy for large projects)
			// fmt.Printf("%s[OK]%s %s\n", colorGreen, colorReset, res.Original.URL)
		} else {
			broken = append(broken, res)
			errStr := ""
			if res.Error != nil {
				errStr = fmt.Sprintf(" (%v)", res.Error)
			}
			fmt.Printf("%s[DEAD: %d]%s %s %s(src: %s:%d)%s\n", 
				colorRed, res.Status, colorReset, 
				res.Original.URL, 
				colorYellow, res.Original.FilePath, res.Original.LineNum, colorReset)
			if errStr != "" {
				fmt.Printf("   -> %s%s%s\n", colorRed, errStr, colorReset)
			}
		}
	}

	fmt.Printf("\n%s--- Summary ---%s\n", colorBold, colorReset)
	fmt.Printf("Total Checked: %d\n", len(results))
	fmt.Printf("Successful:    %s%d%s\n", colorGreen, successCount, colorReset)
	fmt.Printf("Broken:        %s%d%s\n", colorRed, len(broken), colorReset)
	fmt.Printf("Time Taken:    %s%v%s\n", colorCyan, duration, colorReset)

	return len(broken) > 0
}
