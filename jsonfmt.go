package main

import (
	"bytes"
	"encoding/json"
	"flag"
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
