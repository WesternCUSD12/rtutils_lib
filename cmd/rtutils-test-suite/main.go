package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"rtutils_lib/internal/rtconfig"
	"rtutils_lib/internal/testrunner"
)

func main() {
	// Define command-line flags
	reportFlag := flag.String("report", "test-report.json", "Path to output JSON report")
	debugFlag := flag.Bool("debug", false, "Enable debug logging")
	testCasesFlag := flag.String("test-cases", "specs/005-integration-test-suite/contracts/test-cases.json", "Path to test cases JSON file")
	filterServiceFlag := flag.String("service", "", "Filter tests by service (e.g., TicketService)")
	filterValidatorFlag := flag.String("validator", "", "Filter tests by validator (e.g., TicketValidator)")

	flag.Parse()

	// Load configuration from environment variables
	config, err := rtconfig.LoadFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid configuration: %v\n", err)
		os.Exit(1)
	}

	if *debugFlag {
		fmt.Printf("DEBUG: Configuration loaded successfully\n")
		fmt.Printf("DEBUG: RT URL: %s\n", config.URL)
		fmt.Printf("DEBUG: Timeout: %d seconds\n", config.Timeout)
	}

	// Load test cases from JSON file
	testCasesPath := *testCasesFlag
	if !filepath.IsAbs(testCasesPath) {
		// Try to make path relative to current working directory
		cwd, err := os.Getwd()
		if err == nil {
			testCasesPath = filepath.Join(cwd, testCasesPath)
		}
	}

	if *debugFlag {
		fmt.Printf("DEBUG: Loading test cases from: %s\n", testCasesPath)
	}

	loader := testrunner.NewTestCaseLoader(testCasesPath)
	testCases, err := loader.LoadTestCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load test cases: %v\n", err)
		os.Exit(1)
	}

	if *debugFlag {
		fmt.Printf("DEBUG: Loaded %d test cases\n", len(testCases))
	}

	// Apply filters if specified
	if *filterServiceFlag != "" {
		testCases = testrunner.GetTestCasesByService(testCases, *filterServiceFlag)
		if *debugFlag {
			fmt.Printf("DEBUG: Filtered to %d test cases for service: %s\n", len(testCases), *filterServiceFlag)
		}
	}

	if *filterValidatorFlag != "" {
		testCases = testrunner.GetTestCasesByValidator(testCases, *filterValidatorFlag)
		if *debugFlag {
			fmt.Printf("DEBUG: Filtered to %d test cases for validator: %s\n", len(testCases), *filterValidatorFlag)
		}
	}

	// Create and run test runner
	runner := testrunner.NewTestRunner(config, testCases)

	if *debugFlag {
		fmt.Printf("DEBUG: Starting test execution...\n")
	}

	report, err := runner.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: test execution failed: %v\n", err)
		os.Exit(1)
	}

	// Format and output results
	formatter := testrunner.NewReportFormatter(report)

	// Print human-readable report to console
	fmt.Println()
	fmt.Println(formatter.FormatText())

	// Write JSON report to file
	if err := writeReportToFile(*reportFlag, formatter); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write report file: %v\n", err)
	} else if *debugFlag {
		fmt.Printf("DEBUG: Report written to: %s\n", *reportFlag)
	}

	// Exit with appropriate code
	if report.Execution.FailedTests > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

// writeReportToFile writes the test report to a JSON file
func writeReportToFile(filePath string, formatter *testrunner.ReportFormatter) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()

	return formatter.WriteJSON(file, true)
}