package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
)

// Version information - set during build
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Define command line flags
	var inputFile = flag.String("i", "", "Input file path (required)")
	var outputFile = flag.String("o", "index.html", "Output HTML file path (default: index.html)")
	var outputFileLong = flag.String("output-html-path", "index.html", "Output HTML file path (default: index.html)")
	var showVersion = flag.Bool("v", false, "Show version information")
	var showHelp = flag.Bool("h", false, "Show help information")

	// Parse command line flags
	flag.Parse()

	// Handle version flag
	if *showVersion {
		showVersionInfo()
		return
	}

	// Handle help flag
	if *showHelp {
		showHelpInfo()
		return
	}

	// Validate input
	if err := validateInput(*inputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		showUsage()
		os.Exit(1)
	}

	// Determine output file (prefer -o over --output-html-path if both are set)
	finalOutputFile := *outputFile
	if *outputFileLong != "index.html" && *outputFile == "index.html" {
		finalOutputFile = *outputFileLong
	}

	// Display input and output files
	fmt.Printf("Input file: %s\n", *inputFile)
	fmt.Printf("Output file: %s\n", finalOutputFile)

	// Process the files
	if err := processPlanFile(*inputFile, finalOutputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing plan file: %v\n", err)
		os.Exit(1)
	}
}

func validateInput(inputFile string) error {
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file '%s' does not exist", inputFile)
	}

	// Check if input file is readable
	if _, err := os.Open(inputFile); err != nil {
		return fmt.Errorf("cannot read input file '%s': %v", inputFile, err)
	}

	return nil
}

func processPlanFile(inputFile, outputFile string) error {
	fmt.Println("\nProcessing files:")

	// Display file information
	fmt.Printf("Input file: %s\n", inputFile)
	fmt.Printf("Output file: %s\n", outputFile)

	// Read JSON file
	jsonData, err := readJSONFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading JSON file: %v", err)
	}

	// Parse JSON
	var planData interface{}
	if err := json.Unmarshal(jsonData, &planData); err != nil {
		return fmt.Errorf("parsing plan JSON: %v", err)
	}

	fmt.Println("Successfully parsed JSON file!")
	fmt.Printf("JSON contains %d bytes of data\n", len(jsonData))

	// Generate HTML from the parsed plan data
	htmlContent := generateHtml(planData)
	fmt.Printf("Generated HTML content (%d characters)\n", len(htmlContent))

	// Write HTML to output file
	if err := writeHtmlFile(outputFile, htmlContent); err != nil {
		return fmt.Errorf("writing HTML file: %v", err)
	}

	fmt.Printf("Successfully wrote HTML to: %s\n", outputFile)
	fmt.Println("\nFile processing completed!")
	return nil
}

func writeHtmlFile(filePath, htmlContent string) error {
	err := os.WriteFile(filePath, []byte(htmlContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write HTML file %s: %v", filePath, err)
	}
	return nil
}

func readJSONFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}
	return data, nil
}

func showVersionInfo() {
	fmt.Printf("Terraform Plan Visualizer %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func showHelpInfo() {
	fmt.Println("Terraform Plan Visualizer")
	fmt.Println("A tool to convert Terraform plan JSON files into interactive HTML visualizations")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  terraform-plan-visualizer -i <input-file> [-o <output-file>]")
	fmt.Println("  terraform-plan-visualizer -i <input-file> [--output-html-path <output-file>]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -i, -input string        Input Terraform plan JSON file (required)")
	fmt.Println("  -o, -output string       Output HTML file path (default: index.html)")
	fmt.Println("  --output-html-path string")
	fmt.Println("                           Output HTML file path (alternative to -o)")
	fmt.Println("  -v, -version             Show version information")
	fmt.Println("  -h, -help                Show this help information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  terraform-plan-visualizer -i plan.json")
	fmt.Println("  terraform-plan-visualizer -i plan.json -o visualization.html")
	fmt.Println("  terraform-plan-visualizer -i plan.json --output-html-path my-plan.html")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/cloudvic-org/terraform-plan-visualizer")
}

func showUsage() {
	fmt.Println("Usage: terraform-plan-visualizer -i <input-file> [-o <output-file>]")
	fmt.Println("       terraform-plan-visualizer -i <input-file> [--output-html-path <output-file>]")
	fmt.Println("Use -h for more help information")
}
