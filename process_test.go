package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dancsecs/szTest"
)

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

func getProcessedFiles(dir, fName string) (string, string, []string, []string, error) {
	var targetPath string
	var sourceDir string
	var err error
	var gotBytes []byte
	var wntBytes []byte

	targetPath, err = filepath.Abs(filepath.Join(dir, fName))
	if err == nil {
		gotBytes, err = os.ReadFile(targetPath)
	}
	if err == nil {
		sourceDir, err = filepath.Abs(".")
	}
	if err == nil {
		wntBytes, err = os.ReadFile(fName)
	}
	if err != nil {
		return "", "", nil, nil, err
	}
	return targetPath, sourceDir,
		strings.Split(string(gotBytes), "\n"),
		strings.Split(string(wntBytes), "\n"),
		nil
}

func setupDirs(chk *szTest.Chk) {
	chk.T().Helper()
	origOutputDir := outputDir
	origCWD, err := os.Getwd()

	if chk.NoErr(err) {
		outputDir = chk.CreateTmpDir()
		chk.PushPostReleaseFunc(func() error {
			outputDir = origOutputDir
			return os.Chdir(origCWD)
		})
	}
}

func setupTest(
	chk *szTest.Chk, tCleanOnly, tReplace, tForceOverwrite, tVerbose bool,
) {
	chk.T().Helper()
	origOutputDir := outputDir
	origCWD, err := os.Getwd()
	origCleanOnly := cleanOnly
	origReplace := replace
	origForceOverwrite := forceOverwrite
	origVerbose := verbose

	cleanOnly = tCleanOnly
	replace = tReplace
	forceOverwrite = tForceOverwrite
	verbose = tVerbose

	if chk.NoErr(err) {
		outputDir = chk.CreateTmpDir()
		chk.PushPostReleaseFunc(func() error {
			outputDir = origOutputDir
			cleanOnly = origCleanOnly
			replace = origReplace
			forceOverwrite = origForceOverwrite
			verbose = origVerbose
			return os.Chdir(origCWD)
		})
	}
}

// +-------------------------------------------------------+
// | Flag possibilities for type of test.                  |
// +------------+-----------+------------------+-----------+
// | cleanOnly  |  replace  |  forceOverwrite  |  verbose  |
// +------------+-----------+------------------+-----------+
// |  false     |   false   |     false        |   false   |
// |  false     |   true    |     false        |   false   |
// |  true      |   false   |     false        |   false   |
// +------------+-----------+------------------+-----------+
// |  false     |   false   |     false        |   true    |
// |  false     |   true    |     false        |   true    |
// |  true      |   false   |     false        |   true    |
// +------------+-----------+------------------+-----------+
// |  false     |   false   |     true         |   false   |
// |  false     |   true    |     true         |   false   |
// |  true      |   false   |     true         |   false   |
// +------------+-----------+------------------+-----------+
// |  false     |   false   |     true         |   true    |
// |  false     |   true    |     true         |   true    |
// |  true      |   false   |     true         |   true    |
// +------------+-----------+------------------+-----------+.
