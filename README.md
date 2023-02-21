# go-du
`go-du is` a command-line utility tool written in Go programming language that displays disk usage information of a directory or file(s). It calculates the total size of the specified directories or files, and if -s flag is passed, it summarizes the results showing only the total size of the argument(s) specified.

## Usage
The general syntax for using `go-du` is:

`go-du [options] [directory/file(s)]`

The following are the available options:

`-h`: displays the sizes in a human-readable format, such as 1K, 234M, or 2G.

`-b`: sets the block size in bytes. The default value is 1 byte.

`-s`: displays only the total for each argument.

If no directory or file is specified, the current directory is used as the default directory.

## Examples
Display the size of the current directory in bytes:

`go-du`

Display the size of the specified directory in human-readable format:

`go-du -h /path/to/directory`

Display the total size of multiple directories:

`go-du /path/to/directory1 /path/to/directory2`

Display only the total size of the specified directory:

`go-du -s /path/to/directory`

## License
go-du is licensed under the Apache License 2.0. See the LICENSE file for more information.

## Author
go-du was created by Andrew Kuehne in 2023.