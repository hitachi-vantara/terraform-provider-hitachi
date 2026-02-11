package main

import (
	"fmt"
	"os"
	"strings"

	"terraform-provider-hitachi/hitachi/common/config"
)

func dedent(text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, " \t")
	}
	return strings.Join(lines, "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: create_consent_spec_text <output_file>")
		os.Exit(1)
	}

	outputFile := os.Args[1]

	consent := dedent(config.DEFAULT_CONSENT_MESSAGE)
	run := strings.TrimSpace(config.RUN_CONSENT_MESSAGE)

	combined := consent + " " + run

	if err := os.WriteFile(outputFile, []byte(combined), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to %s: %v\n", outputFile, err)
		os.Exit(1)
	}
}
