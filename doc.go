// Package formats json files. It uses tabs for indentation.
//
// Usage
//
// Without an explicit path, it processes the standard input. Given a file, it operates on that file; given a directory, it operates on all .json files in that directory, recursively. (Files starting with a period are ignored.)
// 	jsonfmt -w .
//
// Flags
//
// List:
//
// 	-w 	Do not print reformatted sources to standard output.
//		If a file's formatting is different from jsonfmt's, overwrite it
//		with jsonfmt's version.
package main
