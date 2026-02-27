File shifter - Hexegic Coding Exercise 
---

Usage: rotate [left|right] file_in file_out

Rotates a file by 1 bit, left or right. The file to be rotated is file_in the output file will be created in file_out. If file_out already exists it will be overwritten.

**WARNING** file_in and file_out must be different files, rotation is a destructive operation and the output file is truncated, as such you may loose your data if file_in and file_out are the same file.

The file [approach.md](approach.md) describes the approach I took to creating this tool.

## How to Build
Ensure you have an up to date go toolchain installed and then run:

go build -o renovate .

The compiled binary has been tested on MacOs, and should since it uses only go standard library functions (plus some testing dependencies) it should compile fine for windows or linux.

## Running the tests
Run the following to run the tests:

go test --race ./...

(the race checking should not be required due to there being no parallelism in the implementation, but old habits die hard).