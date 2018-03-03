package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/scanner"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	overwrite bool
)

func init() {
	flag.BoolVar(&overwrite, "w", false, "write result to (source) file instead of stdout")
}

func main() {
	flag.Parse()
	paths := flag.Args()

	if len(paths) == 0 {
		if err := processFile("<standard input>", os.Stdin, os.Stdout, true); err != nil {
			report(err)
		}
		return
	}

	for _, path := range paths {
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			if err := processDir(path, os.Stdout); err != nil {
				report(err)
			}
		default:
			if err := processFile(path, nil, os.Stdout, false); err != nil {
				report(err)
			}
		}
	}
}

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	os.Exit(2)
}

func processDir(path string, out io.Writer) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") && strings.HasSuffix(info.Name(), ".json") {
			return processFile(path, nil, out, false)
		}

		return nil
	})
}

func processFile(filename string, in io.Reader, out io.Writer, stdin bool) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	var tmp interface{}
	if err := json.Unmarshal(src, &tmp); err != nil {

		//attempt to get the line and character offset that caused the error
		if jsonError, ok := err.(*json.SyntaxError); ok {
			line, character, lcErr := lineAndCharacter(string(src), int(jsonError.Offset))
			fmt.Fprintf(os.Stderr, "failed with error: Cannot parse JSON schema due to a syntax error at line %d, character %d: %v\n", line, character, jsonError.Error())
			if lcErr != nil {
				fmt.Fprintf(os.Stderr, "Couldn't find the line and character position of the error due to error %v\n", lcErr)
			}
		}
		if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
			line, character, lcErr := lineAndCharacter(string(src), int(jsonError.Offset))
			fmt.Fprintf(os.Stderr, "failed with error: The JSON type '%v' cannot be converted into the Go '%v' type on struct '%s', field '%v'. See input file line %d, character %d\n", jsonError.Value, jsonError.Type.Name(), jsonError.Struct, jsonError.Field, line, character)
			if lcErr != nil {
				fmt.Fprintf(os.Stderr, "failed with error: Couldn't find the line and character position of the error due to error %v\n", lcErr)
			}
		}

		return err
	}

	res, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		return err
	}

	if !bytes.Equal(src, res) {
		if overwrite {
			return ioutil.WriteFile(filename, res, 0)
		}

		_, err = out.Write(res)
	}
	return err
}

// Source: https://adrianhesketh.com/2017/03/18/getting-line-and-character-positions-from-gos-json-unmarshal-errors/
func lineAndCharacter(input string, offset int) (line int, character int, err error) {
	lf := rune(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0, fmt.Errorf("Couldn't find offset %d within the input.", offset)
	}

	// Humans tend to count from 1.
	line = 1

	for i, b := range input {
		if b == lf {
			line++
			character = 0
		}
		character++
		if i == offset {
			break
		}
	}

	return line, character, nil
}
