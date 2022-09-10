package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/cretz/interflect/genreflect/genreflect"
)

var (
	include regexesFlag
	exclude regexesFlag
	outFile string
	env     stringsFlag
)

func main() {
	// Parse flags
	flag.Var(&include, "include", "Regex to include. If none, all included. Checked for exclusions before exclude.")
	flag.Var(&exclude, "exclude", "Regex to exclude. Checked after include.")
	flag.StringVar(&outFile, "out", "", "File to write. Package will be directory name.")
	flag.Var(&env, "env", "Additional environment variables, like GOOS and GOARCH.")
	flag.Parse()
	if outFile == "" {
		log.Fatal("Missing out file")
	} else if flag.NArg() == 0 {
		log.Fatal("Missing patterns")
	}

	// Get absolute path to out file
	outFile, err := filepath.Abs(outFile)
	if err != nil {
		log.Fatalf("Cannot get absolute path to out file: %v", err)
	}

	// Append env if any present
	if len(env) > 0 {
		env = append(os.Environ(), env...)
	}

	// Generate code
	code, err := genreflect.GenerateReflector(genreflect.GenerateReflectorConfig{
		OutPackage: filepath.Base(filepath.Dir(outFile)),
		Include:    include,
		Exclude:    exclude,
		Patterns:   flag.Args(),
		Env:        env,
	})
	if err != nil {
		log.Fatalf("Failed generating code: %v", err)
	}

	// Write
	if err = os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
		log.Fatalf("Failed creating parent dir: %v", err)
	} else if err = os.WriteFile(outFile, code, 0644); err != nil {
		log.Fatalf("Failed writing out dir: %v", err)
	}
	log.Printf("Wrote %v successfully", outFile)
}

type regexesFlag []*regexp.Regexp

func (r *regexesFlag) String() string { return fmt.Sprintf("%v", *r) }

func (r *regexesFlag) Set(value string) error {
	regex, err := regexp.Compile(value)
	if err == nil {
		*r = append(*r, regex)
	}
	return err
}

type stringsFlag []string

func (s *stringsFlag) String() string { return fmt.Sprintf("%v", *s) }

func (s *stringsFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}
