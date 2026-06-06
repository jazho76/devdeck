package sysreq

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jazho76/devdeck/cli/internal/run"
)

var versionPattern = regexp.MustCompile(`\d+(\.\d+)*[a-z]?`)

func currentVersion(path string, args []string) (string, error) {
	out, err := run.Output(path, args...)
	if err != nil {
		return "", err
	}
	return versionPattern.FindString(out), nil
}

func meetsMinimum(current, min string) bool {
	curNums, curSuffix := parseVersion(current)
	minNums, minSuffix := parseVersion(min)

	for i := 0; i < len(curNums) || i < len(minNums); i++ {
		var c, m int
		if i < len(curNums) {
			c = curNums[i]
		}
		if i < len(minNums) {
			m = minNums[i]
		}
		if c != m {
			return c > m
		}
	}
	return curSuffix >= minSuffix
}

func parseVersion(v string) ([]int, string) {
	var suffix string
	if n := len(v); n > 0 {
		if last := v[n-1]; last >= 'a' && last <= 'z' {
			suffix = string(last)
			v = v[:n-1]
		}
	}
	var nums []int
	for _, part := range strings.Split(v, ".") {
		n, err := strconv.Atoi(part)
		if err != nil {
			break
		}
		nums = append(nums, n)
	}
	return nums, suffix
}
