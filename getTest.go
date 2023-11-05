package main

/*
   Golang To Github Markdown: goToMD.
   Copyright (C) 2023  Leslie Dancsecs

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// "--- PASS: Test_PASS_SampleGoProject (0.0s)".
// "--- FAIL: Test_FAIL_SampleGoProject (0.0s)".
var squashTestTime = regexp.MustCompile(
	`(?m)^(--- (?:PASS|FAIL): .+?) \(\d+\.\d+s\)$`,
)

// Squash runtimes for all tests.
// "FAIL\tgithub.com/dancsecs/szTestToMarkdown/sampleGoProject\t0.0s".
var squashAllTestTime = regexp.MustCompile(`(?m)^FAIL\s(.+?)\s\d+\.\d+s$`)

// Squash (cached) indicator on tests.
// "ok  	example	(cached) coverage: 100.0% of statements".
// "ok  	example	0.269s	coverage: 100.0% of statements".

var squashCached = regexp.MustCompile(
	`(?m)^(ok\s+.+?\s+)(?:\(cached\)|\d+\.\d+s)\s+(.*)$`,
)

func runTest(dir, tests string) (string, string, error) {
	var rawRes []byte
	var args []string
	res := ""

	stat, err := os.Stat(dir)
	if err == nil && !stat.IsDir() {
		err = errors.New("not a directory")
	}
	if err == nil {
		args = []string{"test", "-v", "-cover"}

		if tests != "package" {
			args = append(args, "-run", tests)
		}
		args = append(args, dir)
		c := exec.Command("go", args...) //nolint:gosec // Ok.
		//	c.Dir = dir
		rawRes, _ = c.CombinedOutput() // We expect a general task error.
		if bytes.HasPrefix(rawRes, []byte("testing: warning: no tests to run")) {
			err = errors.New("no tests to run")
		}
	}
	if err == nil {
		res, err = marksToMarkdownHTML(string(rawRes))
	}
	if err == nil {
		res = squashTestTime.ReplaceAllString(res, `${1} (0.0s)`)

		res = squashAllTestTime.ReplaceAllString(res, `FAIL ${1} 0.0s`)

		res = squashCached.ReplaceAllString(res, `${1}${2}`)

		res = strings.ReplaceAll(res, " ", `\unicode{160}`)

		res = strings.ReplaceAll(res, "_", `\unicode{95}`)

		latexRes := ""
		lines := strings.Split(res, "\n")
		for _, line := range lines[:len(lines)-1] {
			latexRes += "$\\small{\\texttt{" + line + "}}$\n<br>\n"
		}
		res = latexRes
	}

	if err == nil {
		return "go " + strings.Join(args, " "), res, nil
	}
	return "", "", err
}

func expandGoTst(cmd string) (string, error) {
	var res string
	var tstCmd string
	dir, action, err := parseCmd(cmd)
	if err == nil {
		tstCmd, res, err = runTest(dir, action)
	}
	if err == nil {
		return "" +
				szTestBgnPrefix + szTstPrefix + cmd + " -->\n" +
				"```bash\n" + tstCmd + "\n```\n\n" +
				res +
				szTestEndPrefix + szTstPrefix + cmd + " -->\n",
			nil
	}
	return "", err
}
