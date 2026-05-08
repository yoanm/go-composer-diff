package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var threshold float64 // Acceptable regression (10% for instance)
	var err error

	if len(os.Args) != 2 {
		log.Fatalf("Missing threshold argument. Usage: %s [threshold_percentage]\n", path.Base(os.Args[0]))
	} else {
		if threshold, err = strconv.ParseFloat(os.Args[1], 64); err != nil {
			log.Fatalln("Threshold must be a valid float")
		}
		if threshold > 100 || threshold <= 0 {
			log.Fatalln("Threshold must be between 1% and 99%")
		}
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 { // Data must come from pipe !
		log.Fatalf(
			"No input detected. Please pipe benchstat output into this tool: "+
				"cat benchstat.out | %s [threshold_percentage]",
			path.Base(os.Args[0]),
		)
	}

	scanner := bufio.NewScanner(os.Stdin)
	// Match lines like: BenchmarkAbc-42  230ns  123ns  +90.00%
	deltaRegex := regexp.MustCompile(`([+-]\d+\.?\d*)%`)

	var regList []string // Regressions list
	for scanner.Scan() {
		line := scanner.Text()
		matches := deltaRegex.FindStringSubmatch(line)

		if len(matches) > 1 {
			delta, err2 := strconv.ParseFloat(matches[1], 64)
			if err2 != nil {
				slog.Error("Error parsing delta. Skipping line", "line", line)
				continue
			}

			// Positive delta means regression (slower)
			if delta > threshold {
				regList = append(
					regList,
					fmt.Sprintf("  %s (%.2f%% slower)", strings.Fields(line)[0], delta),
				)
			}
		}
	}

	if len(regList) > 0 {
		slog.Error(fmt.Sprintf("Performance regression detected (threshold: %.1f%%):\n", threshold))
		for _, reg := range regList {
			slog.Error(reg)
		}
		os.Exit(1)
	}

	slog.Info("All good 🎉.")
}
