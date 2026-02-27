package main

import (
	"fmt"
	"os"
	"rotate/shift"
	"rotate/shift/file"
	"strings"
)

const (
	usage = "Usage: rotate [left|right] file_in file_out"
)

func main() {
	args := os.Args
	if len(args) != 4 {
		fmt.Println(usage)
		os.Exit(-1)
	}
	direction, err := parseDirection(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Println(usage)
		os.Exit(-1)
	}
	fileIn, err := os.Open(args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Println(usage)
		os.Exit(-1)
	}
	fileOut, err := os.Create(args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Println(usage)
		os.Exit(-1)
	}
	err = file.ShiftFile(fileIn, fileOut, direction)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
}

func parseDirection(dir string) (shift.DIRECTION, error) {
	caseless := strings.ToLower(dir)
	switch caseless {
	case "left":
		return shift.LEFT, nil
	case "right":
		return shift.RIGHT, nil
	default:
		return -1, fmt.Errorf("incorrect direction (%s)", dir)
	}
}
