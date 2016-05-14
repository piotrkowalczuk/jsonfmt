# jsonfmt

JSONfmt formats json files. It uses tabs for indentation.

Without an explicit path, it processes the standard input. Given a file, it operates on that file; given a directory, it operates on all .json files in that directory, recursively. (Files starting with a period are ignored.)

## Usage

By default, jsonfmt prints the reformatted sources to standard output.

```
$ jsonfmt test.json
```

With `-w` flag it overwrites given file or files found in directory.

```
$ jsonfmt -w .
```

Using jsonfmt with standard input is straightforward.

```
$ echo '{"b":[    {"firstName":"John", "lastName":"Doe"},    {"firstName":"Anna", "lastName":"Smith"},    {"firstName":"Peter", "lastName":"Jones"}], "a": [3,2,1]}' | jsonfmt
$ {
	"a": [
		3,
		2,
		1
	],
	"b": [
		{
			"firstName": "John",
			"lastName": "Doe"
		},
		{
			"firstName": "Anna",
			"lastName": "Smith"
		},
		{
			"firstName": "Peter",
			"lastName": "Jones"
		}
	]
}
```
## Flags

* `-w` do not print reformatted sources to standard output. If a file's formatting is different from jsonfmt's, overwrite it with jsonfmt's version.
