// Command scan is the entry point for the scan CLI.
//
// It owns flag parsing, wiring, and exit codes. Analysis lives in internal/;
// this package stays thin so the CLI surface is easy to test and easy to
// change.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/zazzcode/scan/internal/lang"
	"github.com/zazzcode/scan/internal/parse"
	"github.com/zazzcode/scan/internal/walk"
)

// Exit codes. CI depends on these, so they are part of the CLI contract.
const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run is the testable entry point: it takes its arguments and streams and
// returns an exit code instead of calling os.Exit.
func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var (
		list   = flags.Bool("list", false, "list the source files under the target path and exit")
		routes = flags.Bool("routes", false, "emit the Express route execution graph and exit")
		asJSON = flags.Bool("json", false, "emit machine-readable JSON instead of text")
	)
	flags.Usage = func() {
		fmt.Fprintln(stderr, "usage: scan --list [--json] <path>")
		fmt.Fprintln(stderr, "       scan --routes [--json] <path>")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return exitUsage
	}
	if *list == *routes || flags.NArg() != 1 {
		flags.Usage()
		return exitUsage
	}

	if *routes {
		return runRoutes(flags.Arg(0), *asJSON, stdout, stderr)
	}

	result, err := walk.Walk(flags.Arg(0), walk.Options{})
	if err != nil {
		fmt.Fprintf(stderr, "scan: %v\n", err)
		return exitError
	}

	if *asJSON {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(newManifest(result)); err != nil {
			fmt.Fprintf(stderr, "scan: %v\n", err)
			return exitError
		}
		return exitOK
	}

	printManifest(stdout, result)
	return exitOK
}

func runRoutes(root string, asJSON bool, stdout, stderr io.Writer) int {
	scanned, err := parse.Routes(root)
	if err != nil {
		fmt.Fprintf(stderr, "scan: %v\n", err)
		return exitError
	}
	if asJSON {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(scanned); err != nil {
			fmt.Fprintf(stderr, "scan: %v\n", err)
			return exitError
		}
		return exitOK
	}
	for _, route := range scanned.Routes {
		fmt.Fprintf(stdout, "%s %s %s:%d\n", route.Method, route.PathPattern, route.Registration.File, route.Registration.StartLine)
	}
	for _, gap := range scanned.Gaps {
		fmt.Fprintf(stdout, "gap %s %s:%d %s\n", gap.Kind, gap.At.File, gap.At.StartLine, gap.Detail)
	}
	return exitOK
}

// manifest is the machine-readable form of a walk result. Field names are part
// of the output contract.
type manifest struct {
	Root        string         `json:"root"`
	Examined    int            `json:"examined"`
	Unsupported int            `json:"unsupported"`
	Supported   int            `json:"supported"`
	ByLanguage  map[string]int `json:"by_language"`
	Files       []manifestFile `json:"files"`
}

type manifestFile struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Size     int64  `json:"size"`
}

func newManifest(result *walk.Result) manifest {
	byLanguage := make(map[string]int, len(result.ByLanguage))
	for language, count := range result.ByLanguage {
		byLanguage[string(language)] = count
	}

	files := make([]manifestFile, 0, len(result.Files))
	for _, file := range result.Files {
		files = append(files, manifestFile{
			Path:     file.Path,
			Language: string(file.Language),
			Size:     file.Size,
		})
	}

	return manifest{
		Root:        result.Root,
		Examined:    result.Examined,
		Unsupported: result.Unsupported,
		Supported:   len(result.Files),
		ByLanguage:  byLanguage,
		Files:       files,
	}
}

func printManifest(w io.Writer, result *walk.Result) {
	fmt.Fprintf(w, "root:        %s\n", result.Root)
	fmt.Fprintf(w, "examined:    %d\n", result.Examined)
	fmt.Fprintf(w, "supported:   %d\n", len(result.Files))
	fmt.Fprintf(w, "unsupported: %d\n", result.Unsupported)
	fmt.Fprintln(w)

	for _, language := range sortedLanguages(result.ByLanguage) {
		fmt.Fprintf(w, "%-12s %d\n", language, result.ByLanguage[language])
	}

	fmt.Fprintln(w)
	for _, file := range result.Files {
		fmt.Fprintln(w, file.Path)
	}
}

// sortedLanguages returns the languages present in counts in a stable order.
// Map iteration order is deliberately never used for output.
func sortedLanguages(counts map[lang.Language]int) []lang.Language {
	languages := make([]lang.Language, 0, len(counts))
	for language := range counts {
		languages = append(languages, language)
	}
	sort.Slice(languages, func(i, j int) bool { return languages[i] < languages[j] })
	return languages
}
