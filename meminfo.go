package main

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
	"regexp"
	"strconv"
)

var (
	pageSizeRegex   *regexp.Regexp = regexp.MustCompile(`page size of (\d*) bytes`)
	pageNumberRegex *regexp.Regexp = regexp.MustCompile(`([a-zA-Z "-]*): *([\d]*)\.`)
)

func meminfo() (MemInfo, error) {
	m := MemInfo{
		values: map[string]int64{},
	}

	cmd := exec.Command("vm_stat")
	out, err := cmd.Output()
	if err != nil {
		return m, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))

	if !scanner.Scan() {
		return m, io.ErrUnexpectedEOF
	}
	matches := pageSizeRegex.FindStringSubmatch(scanner.Text())

	pageSize, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return m, err
	}

	m.pageSize = pageSize

	for scanner.Scan() {
		matches := pageNumberRegex.FindStringSubmatch(scanner.Text())
		noOfPages, err := strconv.ParseInt(matches[2], 10, 64)
		if err != nil {
			return m, err
		}
		m.values[matches[1]] = noOfPages
	}

	return m, nil
}

type MemInfo struct {
	pageSize int64
	values   map[string]int64
}
